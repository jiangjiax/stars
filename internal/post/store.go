package post

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Store 管理文章存储
type Store struct {
	mu sync.RWMutex

	// 主存储
	posts map[string]*Post // slug -> post

	// 元数据缓存
	metas map[string]*PostMeta // slug -> meta

	// 索引
	dateIndex   []*PostMeta            // 按日期排序的文章元数据
	seriesIndex map[string][]*PostMeta // 系列名 -> 该系列的文章元数据
	tagsIndex   map[string][]*PostMeta // 标签 -> 带此标签的文章元数据

	// 新增：筛选项统计缓存
	seriesStats map[string]int // 系列名 -> 文章数量
	tagsStats   map[string]int // 标签名 -> 文章数量
	allTags     []string       // 所有标签（已排序）
}

// New 创建新的文章存储
func New() *Store {
	return &Store{
		posts:       make(map[string]*Post),
		metas:       make(map[string]*PostMeta),
		seriesIndex: make(map[string][]*PostMeta),
		tagsIndex:   make(map[string][]*PostMeta),
		seriesStats: make(map[string]int),
		tagsStats:   make(map[string]int),
	}
}

// Add 添加文章
func (s *Store) Add(post *Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if post.Slug == "" {
		return fmt.Errorf("post slug cannot be empty")
	}

	// 计算阅读时间（如果还没有计算）
	if post.ReadingTime == 0 {
		post.ReadingTime = calculateReadingTime(post.RawContent)
	}

	// 更新主存储
	s.posts[post.Slug] = post

	// 创建元数据副本
	meta := &PostMeta{
		Title:       post.Title,
		Date:        post.Date,
		Description: post.Description,
		Tags:        post.Tags,
		Slug:        post.Slug,
		Series:      post.Series,
		SeriesOrder: post.SeriesOrder,
		ReadingTime: post.ReadingTime, // 确保复制阅读时间
	}
	s.metas[post.Slug] = meta

	// 更新索引
	s.updateDateIndex(post.Slug, meta)
	s.updateSeriesIndex(post.Slug, meta)
	s.updateTagsIndex(post.Slug, meta)

	// 更新系列统计
	if post.Series != "" {
		s.seriesStats[post.Series]++
	}

	// 更新标签统计
	for _, tag := range post.Tags {
		s.tagsStats[tag]++
	}

	// 更新标签列表（如果有新标签）
	s.updateTagsList()

	return nil
}

// Get 获取文章
func (s *Store) Get(slug string) (*Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[slug]
	if !ok {
		return nil, fmt.Errorf("post not found: %s", slug)
	}

	return post, nil
}

// List 获取所有文章
func (s *Store) List() []*Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	// 按日期排序（最新的在前）
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts
}

// LoadFromFS 从文件系统加载文章
func (s *Store) LoadFromFS(fsys fs.FS, dir string) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := fs.ReadFile(fsys, filepath.Join(dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}

		post, err := ParseContent(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse post %s: %w", entry.Name(), err)
		}

		if err := s.Add(post); err != nil {
			return fmt.Errorf("failed to add post %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// GetAll 获取所有文章（按日期排序）
func (s *Store) GetAll() []*Post {
	return s.List() // 直接使用已有的 List 方法
}

// GetBySlug 通过 slug 获取文章
func (s *Store) GetBySlug(slug string) (*Post, error) {
	return s.Get(slug) // 直接使用已有的 Get 方法
}

// ListPaged 获取分页的文章列表
func (s *Store) ListPaged(page, pageSize int) ([]*Post, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 使用 dateIndex 获取已排序的文章
	total := len(s.dateIndex)
	start := (page - 1) * pageSize
	if start >= total {
		return nil, total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	// 从 dateIndex 获取对应的文章
	posts := make([]*Post, 0, end-start)
	for i := start; i < end; i++ {
		meta := s.dateIndex[i]
		if post, ok := s.posts[meta.Slug]; ok {
			posts = append(posts, post)
		}
	}

	return posts, total
}

// ListByTag 获取标签下的分页文章列表
func (s *Store) ListByTag(tag string, page, pageSize int) ([]*Post, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 获取标签下的所有文章元数据
	tagMetas := s.tagsIndex[tag]
	total := len(tagMetas)

	start := (page - 1) * pageSize
	if start >= total {
		return nil, total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	// 将 PostMeta 转换为 Post
	posts := make([]*Post, 0, end-start)
	for i := start; i < end; i++ {
		meta := tagMetas[i]
		if post, ok := s.posts[meta.Slug]; ok {
			posts = append(posts, post)
		}
	}

	return posts, total
}

// 获取系列文章
func (s *Store) GetSeriesPosts(series string) []*PostMeta {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.seriesIndex[series]
}

// 更新日期索引
func (s *Store) updateDateIndex(slug string, meta *PostMeta) {
	// 如果是新文章，直接添加到日期索引
	if len(s.dateIndex) == 0 {
		s.dateIndex = append(s.dateIndex, meta)
		return
	}

	// 按日期插入到正确的位置（保持降序）
	inserted := false
	for i, p := range s.dateIndex {
		if meta.Date.After(p.Date) {
			// 在此位置插入
			s.dateIndex = append(s.dateIndex[:i], append([]*PostMeta{meta}, s.dateIndex[i:]...)...)
			inserted = true
			break
		}
	}

	// 如果是最旧的文章，添加到末尾
	if !inserted {
		s.dateIndex = append(s.dateIndex, meta)
	}
}

// 更新系列索引
func (s *Store) updateSeriesIndex(slug string, meta *PostMeta) {
	if meta.Series == "" {
		return
	}

	// 获取当前系列的文章列表
	posts := s.seriesIndex[meta.Series]

	// 移除旧的索引（如果存在）
	for i, p := range posts {
		if p.Slug == slug {
			posts = append(posts[:i], posts[i+1:]...)
			break
		}
	}

	// 添加新的索引
	posts = append(posts, meta)

	// 按 SeriesOrder 排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].SeriesOrder < posts[j].SeriesOrder
	})

	s.seriesIndex[meta.Series] = posts
}

// 更新标签索引
func (s *Store) updateTagsIndex(slug string, meta *PostMeta) {
	// 先从所有标签索引中移除这篇文章
	for tag, posts := range s.tagsIndex {
		for i, p := range posts {
			if p.Slug == slug {
				s.tagsIndex[tag] = append(posts[:i], posts[i+1:]...)
				break
			}
		}
	}

	// 添加到新的标签索引
	for _, tag := range meta.Tags {
		if s.tagsIndex[tag] == nil {
			s.tagsIndex[tag] = make([]*PostMeta, 0)
		}
		s.tagsIndex[tag] = append(s.tagsIndex[tag], meta)

		// 按日期排序（最新的在前）
		sort.Slice(s.tagsIndex[tag], func(i, j int) bool {
			return s.tagsIndex[tag][i].Date.After(s.tagsIndex[tag][j].Date)
		})
	}
}

// 获取标签下的文章
func (s *Store) GetTagPosts(tag string) []*PostMeta {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tagsIndex[tag]
}

// 获取所有标签及其文章数量
func (s *Store) GetTagsCount() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counts := make(map[string]int)
	for tag, posts := range s.tagsIndex {
		counts[tag] = len(posts)
	}
	return counts
}

// 获取所有系列及其文章数量
func (s *Store) GetSeriesCount() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counts := make(map[string]int)
	for series, posts := range s.seriesIndex {
		counts[series] = len(posts)
	}
	return counts
}

// ListBySeries 获取系列下的分页文章列表
func (s *Store) ListBySeries(series string, page, pageSize int) ([]*Post, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seriesPosts := make([]*Post, 0)
	for _, post := range s.posts {
		if post.Series == series {
			seriesPosts = append(seriesPosts, post)
		}
	}

	// 按 SeriesOrder 排序
	sort.Slice(seriesPosts, func(i, j int) bool {
		return seriesPosts[i].SeriesOrder < seriesPosts[j].SeriesOrder
	})

	total := len(seriesPosts)
	start := (page - 1) * pageSize
	if start >= total {
		return nil, total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return seriesPosts[start:end], total
}

// 更新标签列表
func (s *Store) updateTagsList() {
	tags := make([]string, 0, len(s.tagsStats))
	for tag := range s.tagsStats {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	s.allTags = tags
}

// 获取系列统计
func (s *Store) GetSeriesStats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回一个副本以防止并发修改
	stats := make(map[string]int, len(s.seriesStats))
	for k, v := range s.seriesStats {
		stats[k] = v
	}
	return stats
}

// 获取标签统计
func (s *Store) GetTagsStats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]int, len(s.tagsStats))
	for k, v := range s.tagsStats {
		stats[k] = v
	}
	return stats
}

// 获取所有标签（已排序）
func (s *Store) GetAllTags() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tags := make([]string, len(s.allTags))
	copy(tags, s.allTags)
	return tags
}

// ListByTagAndSeries 同时按标签和系列筛选文章
func (s *Store) ListByTagAndSeries(tag, series string, page, pageSize int) ([]*Post, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 获取同时满足条件的文章
	var filteredPosts []*Post
	for _, post := range s.posts {
		if post.Series == series {
			for _, t := range post.Tags {
				if t == tag {
					filteredPosts = append(filteredPosts, post)
					break
				}
			}
		}
	}

	// 按日期排序
	sort.Slice(filteredPosts, func(i, j int) bool {
		return filteredPosts[i].Date.After(filteredPosts[j].Date)
	})

	total := len(filteredPosts)
	start := (page - 1) * pageSize
	if start >= total {
		return nil, total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return filteredPosts[start:end], total
}

// GetPostsByTag 获取指定标签的所有文章
func (s *Store) GetPostsByTag(tag string) []*Post {
	var posts []*Post
	for _, p := range s.posts {
		for _, t := range p.Tags {
			if t == tag {
				posts = append(posts, p)
				break
			}
		}
	}
	return posts
}

// GetPostsBySeries 获取指定系列的所有文章
func (s *Store) GetPostsBySeries(series string) []*Post {
	var posts []*Post
	for _, p := range s.posts {
		if p.Series == series {
			posts = append(posts, p)
		}
	}
	return posts
}

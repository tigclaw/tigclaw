import blogPostsData from './blogPosts.json'

export interface BlogPost {
  slug: string
  title: string
  excerpt: string
  date: string
  readTime: string
  tag: string
  tagColor: string
  content: string
}

export const blogPosts: BlogPost[] = blogPostsData as BlogPost[]
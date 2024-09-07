import type { MetadataRoute } from 'next'
 
export default function sitemap(): MetadataRoute.Sitemap {
  return [
    {
        url: 'https://bluthinator.com',
        lastModified: new Date(),
        priority: 1,
    },
    { 
        url: 'https://bluthinator.com/about',
        lastModified: new Date(),
        priority: 0.3,
    },
    {
        url: 'https://bluthinator.com/random',
        lastModified: new Date(),
        priority: 0.5,
    }
  ]
}
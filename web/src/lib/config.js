// Default branding configuration
export const defaultConfig = {
  title: 'GoShort - URL Shortener',
  description: 'GoShort is a powerful and user-friendly URL shortener. Simplify, manage, and track your links with ease.',
  keywords: 'URL shortener, GoShort, link management, shorten URLs, track links',
  author: 'GoShort Team',
  themeColor: '#4caf50',
  logoText: 'GoShort',
  primaryColor: '#3b82f6', // blue-600
  secondaryColor: '#10b981', // emerald-500
  headerTitle: 'GoShort - URL Shortener',
  footerText: 'View the project on GitHub',
  footerLink: 'https://github.com/kek-Sec/GoShort'
};

// Load config from the server at runtime
export const loadConfig = async () => {
  try {
    const response = await fetch('/api/v1/config');
    if (response.ok) {
      const config = await response.json();
      return { ...defaultConfig, ...config };
    }
    return defaultConfig;
  } catch (error) {
    console.warn('Failed to load custom configuration, using defaults', error);
    return defaultConfig;
  }
};

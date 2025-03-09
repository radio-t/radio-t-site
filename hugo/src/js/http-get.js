const cache = new Map();
const throttleMap = new Map();

/**
 * Fetch wrapper with basic caching and throttling
 * @param {string} url - URL to fetch
 * @param {Object} options - Fetch options
 * @returns {Promise<any>} - Promise with parsed JSON response
 */
async function fetchWithCache(url, options = {}) {
  // Handle params if present
  if (options.params) {
    const urlObj = new URL(url);
    Object.entries(options.params).forEach(([key, value]) => {
      urlObj.searchParams.append(key, value);
    });
    url = urlObj.toString();
    // Remove params from options to avoid duplication
    delete options.params;
  }

  // sort options to make sure the cache key is consistent
  const cacheKey = url + (Object.keys(options).length ? `-${JSON.stringify(Object.keys(options).sort().reduce((acc, key) => ({
    ...acc, [key]: options[key]
  }), {}))}` : '');

  if (throttleMap.has(cacheKey)) {
    return throttleMap.get(cacheKey);
  }

  if (cache.has(cacheKey)) {
    return cache.get(cacheKey);
  }

  const promise = (async () => {
    try {
      const response = await fetch(url, options);

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();
      cache.set(cacheKey, data);
      setTimeout(() => throttleMap.delete(cacheKey), 0);
      return data;
    } catch (error) {
      throttleMap.delete(cacheKey);
      throw error;
    }
  })();

  throttleMap.set(cacheKey, promise);

  return promise;
}

export default {
  get: async (url, options = {}) => {
    const data = await fetchWithCache(url, {
      method: 'GET', ...options
    });
    return {data};
  }
};
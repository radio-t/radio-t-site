const cache = new Map();
const throttleMap = new Map();

/**
 * Fetch JSON from URL with caching and throttling
 * @param {string} url - URL to fetch
 * @returns {Promise<any>} - Promise with parsed JSON response
 */
export async function fetchJSON(url) {
  const cacheKey = url;

  if (throttleMap.has(cacheKey)) {
    return throttleMap.get(cacheKey);
  }

  if (cache.has(cacheKey)) {
    return cache.get(cacheKey);
  }

  const promise = (async () => {
    try {
      const response = await fetch(url, { method: 'GET' });

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
// helper class for finding the closest key in a map using binary search
class ClosestMap {
  /**
   * @param {ClosestMapEntry[]} entries
   */
  constructor(entries = []) {
    this.map = new Map(entries);
    this._sortKeys();
  }

  /**
   * Add a key-value pair to the map
   * @param {number} key
   * @param {*} value
   */
  add(key, value) {
    this.map.set(key, value);
    this._sortKeys();
  }

  _sortKeys() {
    this.sortedKeys = [...this.map.keys()].sort((a, b) => a - b);
  }

  /**
   * Get the value of the closest key less than or equal to the given key
   * @param {number} key
   * @returns {*|null} value or null
   */
  getFloor(key) {
    let left = 0;
    let right = this.sortedKeys.length - 1;
    let result = null;

    while (left <= right) {
      const mid = Math.floor((left + right) / 2);
      if (this.sortedKeys[mid] <= key) {
        result = this.sortedKeys[mid];
        left = mid + 1;
      } else {
        right = mid - 1;
      }
    }

    return this.map.get(result);
  }
}

export default ClosestMap;

// ============== debounce.js =====================

function debounce(callback, delay = 0) {
  let timeoutId;

  return function (...args) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => {
      timeoutId = undefined;
      callback.apply(this, args);
    }, Math.max(0, delay));
  };
}

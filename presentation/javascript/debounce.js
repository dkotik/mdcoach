// ============== debounce.js =====================

function debounce(func, wait) {
  let timeoutId;

  return function (...args) {
    // Clear the previous timer if the function is called again before the delay finishes
    clearTimeout(timeoutId);

    // Set a new timer to execute the function after the specified delay
    timeoutId = setTimeout(() => {
      func.apply(this, args);
    }, wait);
  };
}

// https://xparkmedia.com/blog/enter-fullscreen-mode-javascript/
// edited by dropping moz methods
export const toggleFullScreen = function() {
  if (
    (document.fullScreenElement && document.fullScreenElement !== null) || // alternative standard method
    (!document.mozFullScreen && !document.webkitIsFullScreen)
  ) {
    // current working methods
    if (document.documentElement.requestFullscreen) {
      document.documentElement.requestFullscreen();
    } else if (document.documentElement.requestFullScreen) {
      document.documentElement.requestFullScreen();
    } else if (document.documentElement.webkitRequestFullScreen) {
      document.documentElement.webkitRequestFullScreen(
        Element.ALLOW_KEYBOARD_INPUT
      );
    }
  } else {
    if (document.exitFullscreen) {
      document.cancelFullScreen();
    } else if (document.exitFullscreen) {
      document.cancelFullScreen();
    } else if (document.webkitCancelFullScreen) {
      document.webkitCancelFullScreen();
    }
  }
}

document.addEventListener('keydown', (event) => {
  if (event.code === 'KeyF') {
    toggleFullScreen()
  }
})

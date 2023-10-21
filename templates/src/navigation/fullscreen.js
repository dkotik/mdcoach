// https://xparkmedia.com/blog/enter-fullscreen-mode-javascript/

// TODO: mozRequestFullScreen() is deprecated.

export const toggleFullScreen = function() {
  if (
    (document.fullScreenElement && document.fullScreenElement !== null) || // alternative standard method
    (!document.mozFullScreen && !document.webkitIsFullScreen)
  ) {
    // current working methods
    if (document.documentElement.requestFullScreen) {
      document.documentElement.requestFullScreen();
    } else if (document.documentElement.mozRequestFullScreen) {
      document.documentElement.mozRequestFullScreen();
    } else if (document.documentElement.webkitRequestFullScreen) {
      document.documentElement.webkitRequestFullScreen(
        Element.ALLOW_KEYBOARD_INPUT
      );
    }
  } else {
    if (document.cancelFullScreen) {
      document.cancelFullScreen();
    } else if (document.mozCancelFullScreen) {
      document.mozCancelFullScreen();
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

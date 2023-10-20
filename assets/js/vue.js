import { toggleFullScreen } from "./fullscreen.js";

export var v = new Vue({
  el: "main",
  data: {
    currentSlide: null,
    slides: slides,
    leftToRight: true,
    daggers: [],
    curtain: false,
    selected: [],
    deleted: [],
    navigator: false,
    controls: false
  },
  computed: {
    slide: {
      get: function() {
        return this.currentSlide;
      },
      set: function(value) {
        this.navigator = false;
        this.controls = false;
        // console.log('current slide set: ' + this.currentSlide + ' =>> ' + value);
        if (value < 1) value = 1;
        else if (value > this.slides.length) value = this.slides.length;
        if (this.curtain || value == this.currentSlide) return;
        location.hash = "slide" + value;
        var direction = value > this.currentSlide;
        // if (direction !== this.leftToRight) console.log('direction changed');
        // -1 is because slides and classes vars start index at 0 rather than 1
        this.$emit(
          "slide",
          value - 1,
          this.currentSlide - 1,
          direction !== this.leftToRight
        );
        this.currentSlide = value;
        this.leftToRight = direction;
        document.title = document.title.replace(
          /^([^\.]+\.\s)?(.*)$/,
          value + "/" + (this.slides.length - 1) + ". $2"
        );
      }
    }
  },
  methods: {
    toggleFullScreen: toggleFullScreen,
    toggleContrast: function(e) {
      if (document.body.classList.contains("flipContrast")) {
        delete window.localStorage["flipContrast"];
        document.body.classList.remove("flipContrast");
      } else {
        window.localStorage["flipContrast"] = "true";
        document.body.classList.add("flipContrast");
      }
    },
    hashEvent: function(e) {
      this.slide = parseInt(location.hash.substr(6)) || 1;
    },
    keyboardEvent: function(e) {
      // console.log(e);
      if (e.code == "Escape") {
        this.navigator = false;
        this.curtain = false;
        e.preventDefault();
        return false;
      } else if (e.code == "Period") {
        this.curtain = !this.curtain;
        this.navigator = false;
      } else if (
        e.code == "ArrowRight" ||
        e.code == "ArrowDown" ||
        e.code == "PageDown" ||
        e.code == "Space"
      ) {
        if (e.shiftKey) {
          this.daggers = [];
        } else if (this.daggers.length) {
          // console.log(this.daggers);
          for (var i = 0; i < this.daggers.length; ++i) {
            if (!this.daggers[i].classList.contains("activated")) {
              this.daggers[i].classList.add("activated");
              return;
            }
          }
          this.daggers = [];
        }
        this.slide += e.shiftKey ? 5 : 1;
        // presentationForward();
      } else if (
        e.code == "ArrowLeft" ||
        e.code == "ArrowUp" ||
        e.code == "PageUp"
      ) {
        this.daggers = [];
        this.slide -= e.shiftKey ? 5 : 1;
        // presentationBack();
      }
    }
  },
  mounted: function() {
    document.title = "1. " + document.title;
    if (window.localStorage["flipContrast"])
      document.body.classList.add("flipContrast");
    this.$nextTick(function() {
      window.addEventListener("keydown", this.keyboardEvent);
      window.addEventListener("hashchange", this.hashEvent);
      // setTimeout(this.hashEvent, 3000);
      // this.currentSlide = parseInt(location.hash.substr(6));
      this.hashEvent();
    });
    this.$on("ready", function(slide) {
      try {
        this.daggers = slide.$el.querySelectorAll(":scope > ul > li");
      } catch (e) {
        this.daggers = [];
      }
      // console.log('this.daggers' + (this.currentSlide-1) + '|');
      // console.log(this.daggers);
    });
  }
});

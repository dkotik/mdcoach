Vue.component('slide', {
  props: {
    slideid: { required: true, type: Number }
  },
  template: '<transition appear :name="left2right ? \'slide-right\' : \'slide-left\'" v-on:after-enter="ready()"><section v-if="display" v-html="slides[slideid]"></section></transition>',
  methods: {
    ready: function() {
      this.$root.$emit('ready', this);
      // reset all scaling
      this.fontSize = 100;
      this.scale = 1;
      this.padding = 0;
      this.$el.style.fontSize = '100%';
      this.$el.style.transform = 'scale(1)';
      setTimeout(this.autoScale, 50);
    },
    conceal: function() {
      this.display = false;
    },
    switch: function(slideNumber, previousSlide, directionChanged) {
      if (slideNumber == this.slideid) {
        this.left2right = slideNumber > previousSlide;
        this.display = true;
        window.addEventListener('resize', this.ready);
      } else if (previousSlide == this.slideid) {
        window.removeEventListener('resize', this.ready);
        if(directionChanged) this.left2right = !this.left2right;
        // required for transition to work properly
        this.$nextTick(this.conceal);
      }
    },
    autoScale: function () {
      if (!this.display) return;
      var ratio = this.$root.$el.offsetHeight / (this.$el.offsetHeight + 20);
      // disabled scale storage
      // if (this.scale) {
      //   this.$el.style.transform = 'scale(' + this.scale + ')';
      //   this.$el.style.fontSize = this.fontSize + '%';
      //   this.$el.style.paddingTop = this.padding + 'px';
      //   // console.log('recalled scale' + this.scale);
      //   return
      // } else
      if (ratio > 1) {
        var padding = (this.$root.$el.offsetHeight - this.$el.offsetHeight) / 3;
        if (Math.abs(this.padding - padding) < 1) return; // stop scaling if padding difference is too small
        this.$el.style.paddingTop = padding + 'px';
        this.padding = padding;
        // this.scale = 1;
        // console.log('setting margin');
      // this is for old figure scaling?
      // } else if (this.$el.childNodes.length == 1 && this.$el.firstChild.nodeName == 'FIGURE') {
      //   return false;
      } else if (this.fontSize > 50) {
        this.fontSize *= 0.80;
        this.$el.style.fontSize = this.fontSize + '%';
        // console.log('reduced font: ' + this.fontSize);
      } else {
        if (this.scale === ratio) return; // stop scaling!
        this.$el.style.transform = 'scale(' + ratio + ')';
        this.scale = ratio;
      }
      setTimeout(this.autoScale, 150);
      // console.log(`scaled to`, ratio);
    }
  },
  created: function() { this.$root.$on('slide', this.switch); },
  data: function() { return {
    scale: 1,
    padding: 0,
    fontSize: 100,
    display: false, left2right: true}; }
});

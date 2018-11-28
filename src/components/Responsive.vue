<template>
  <div class="wrapper">
    <slot v-if="width" :width="width" :height="height"/>
  </div>
</template>

<script>
export default {
  props: {
    mode: {
      type: String,
      default: "both" // maybe later I'll need modes based on only width/height
    },
    proportion: {
      type: Number,
      default: 0.625
    }
  },

  data() {
    return {
      wrapperWidth: 0,
      wrapperHeight: 0
    };
  },

  computed: {
    proportionHeight() {
      return this.wrapperWidth * this.proportion;
    },

    width() {
      if (this.proportionHeight > this.wrapperHeight) {
        return this.wrapperHeight / this.proportion;
      }
      return this.wrapperWidth;
    },

    height() {
      if (this.proportionHeight > this.wrapperHeight) {
        return this.wrapperHeight;
      }
      return this.proportionHeight;
    }
  },

  mounted() {
    this.updateSize();
    window.addEventListener("resize", this.updateSize);
  },

  beforeDestroy() {
    window.removeEventListener("resize", this.updateSize);
  },

  methods: {
    updateSize() {
      this.wrapperWidth = this.$el.clientWidth;
      this.wrapperHeight = this.$el.clientHeight;
    }
  }
};
</script>

<style scoped>
.wrapper {
  height: 100%;
}
</style>

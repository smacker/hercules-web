<template>
  <div><slot v-if="width" :width="width" :height="height" /></div>
</template>

<script>
export default {
  props: {
    proportion: {
      type: Number,
      default: 0.625
    }
  },

  data() {
    return {
      width: 0
    };
  },

  computed: {
    height() {
      return this.width * this.proportion;
    }
  },

  mounted() {
    this.updateWidth();
    window.addEventListener('resize', this.updateWidth);
  },

  beforeDestroy() {
    window.removeEventListener('resize', this.updateWidth);
  },

  methods: {
    updateWidth() {
      this.width = this.$el.clientWidth;
    }
  }
};
</script>

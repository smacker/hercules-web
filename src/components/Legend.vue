<template>
  <ul class="legend">
  <li class="item" v-for="(item, i) in items" :key="item">
      <span class="color" :style="{backgroundColor: colors[i]}" /> {{item}}
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    items: {
      type: Array,
      required: true
    },
    colors: {
      type: Array,
      required: true
    }
  }
};
</script>

<style scoped>
.legend {
  margin: 0;
  padding: 0;
  list-style: none;
}
.legend li {
  margin: 0;
  padding: 0;
}
.color {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 1px solid #000;
}
</style>


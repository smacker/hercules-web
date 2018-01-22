<template>
  <div class="error">
    <h3>Oops! There is an error:</h3>
    <p>{{ msg }}</p>
    <p><a href="/">Try another repository</a></p>
  </div>
</template>

<script>
export default {
  props: ['msg']
};
</script>

<style scoped>
.error {
  color: #f5222d;
}
</style>

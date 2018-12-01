<template>
  <div class="error">
    <h3>Oops! There is an error:</h3>
    <p>{{ msg.toString() }}</p>
    <p>
      <a href="/">Try another repository</a>
    </p>
  </div>
</template>

<script>
export default {
  props: ["msg"]
};
</script>

<style scoped>
.error {
  text-align: center;
  color: #f5222d;
}
</style>

<template>
  <div class="loading">
    <div class="loading__spinner">
      <spinner/>
    </div>
    <div class="loading__text">Fetching &amp; calculating....
      <br>Please wait, it can take few seconds.
    </div>
  </div>
</template>

<script>
import Spinner from "@/components/Spinner";

export default {
  components: {
    Spinner
  }
};
</script>

<style scoped>
.loading {
  font-size: 2em;
  text-align: center;
}

.loading__spinner {
  height: 5em;
  padding: 1.5em 0 2.5em;
  margin-bottom: 10px;
}
</style>

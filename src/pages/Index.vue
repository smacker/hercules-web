<template>
  <div class="index-page">
    <main-form :on-submit="submit" />
  </div>
</template>

<script>
import MainForm from '@/components/MainForm';

export default {
  components: {
    MainForm
  },

  methods: {
    submit(repo) {
      this.$router.push({ path: `/${repo}/burndown` });
    }
  }
};
</script>

<style scoped>
.main-form {
  width: 500px;
  margin: 50px auto;
}
</style>

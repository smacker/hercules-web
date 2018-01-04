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

  data() {
    return {
      error: null
    };
  },

  methods: {
    submit(repo) {
      // FIXME add simple validation
      this.$router.push({ path: `/${repo}/burndown` });
    }
  }
};
</script>

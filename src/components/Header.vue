<template>
  <el-menu :default-active="page" mode="horizontal" class="header" @select="handleSelect">
    <el-menu-item index="overall">Overview</el-menu-item>
    <el-menu-item index="people">By people</el-menu-item>
    <el-menu-item index="files">By files</el-menu-item>
    <li class="repo">
      <form @submit.prevent="handleRepoSubmit">
        <el-input v-model="repoVal" class="input" :disabled="loading">
          <el-button slot="append" icon="el-icon-search" native-type="submit" />
        </el-input>
      </form>
    </li>
    <li class="el-menu-item slot">
      <slot />
    </li>
  </el-menu>
</template>

<script>
export default {
  props: ['page', 'repo', 'loading'],

  data() {
    return {
      repoVal: this.repo
    };
  },

  methods: {
    handleSelect(key) {
      if (key === 'overall') {
        this.$router.push({
          name: 'project',
          params: { repo: this.repo }
        });
        return;
      }
      this.$router.push({
        name: key,
        params: { repo: this.repo }
      });
    },

    handleRepoSubmit() {
      this.$router.push({
        name: this.$route.name,
        params: { repo: this.repoVal }
      });
    }
  }
};
</script>

<style scoped>
.header {
  margin-bottom: 20px;
}

.repo {
  float: left;
  padding-left: 70px;
}

.slot {
  float: right;
}

.input {
  width: 400px;
  margin: 10px 0;
}
</style>

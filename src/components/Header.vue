<template>
  <el-menu :default-active="page" mode="horizontal" class="header" @select="handleSelect">
    <el-menu-item index="overall">Overview</el-menu-item>
    <el-menu-item index="people">By people</el-menu-item>
    <el-menu-item index="files">By files</el-menu-item>
    <li class="repo">
      <el-input :value="repo" class="input" :disabled="loading">
        <el-button slot="append" icon="el-icon-search" />
      </el-input>
    </li>
    <li class="el-menu-item slot">
      <slot />
    </li>
  </el-menu>
</template>

<script>
export default {
  props: ['page', 'repo', 'loading'],

  methods: {
    handleSelect(key) {
      if (key === 'overall') {
        this.$router.push({ path: `/${this.repo}/burndown` });
        return;
      }
      this.$router.push({ path: `/${this.repo}/burndown/${key}` });
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

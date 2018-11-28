<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: "app"
};
</script>

<style>
* {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
}

html,
body,
#app,
.page {
  height: 100%;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica,
    Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
}

.page {
  display: flex;
  flex-direction: column;
}

.page > .header {
  flex: 0 0 auto;
}

.page-body {
  flex: 1 1 auto;
  height: 100%;
  margin-top: -81px;
  padding-top: 81px;
  padding-left: 20px;
  padding-right: 20px;
}
</style>

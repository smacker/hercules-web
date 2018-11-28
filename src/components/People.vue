<template>
  <div class="people">
    <div class="list">
      <div class="people-item" @click="onClick(null)">Show all</div>
      <div
        v-for="item in items"
        :key="item.value"
        :class="{'people-item': true, selected: selected && item.value === selected.value}"
        @click="onClick(item)"
      >
        <span class="color" :style="{backgroundColor: item.color}"/>
        <span class="text">{{item.label}}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: ["items", "selected", "onClick"]
};
</script>

<style scoped>
.people {
  height: 100%;
  background: #fff;
  border: 1px solid #eaecef;
  line-height: 20px;
}

.list {
  height: 100%;
  overflow: auto;
}

.people-item {
  padding: 5px 10px;
  border-bottom: 1px solid #eaecef;
  cursor: pointer;
}

.people-item:hover {
  background: #f6f8fa;
}

.people-item.selected {
  background: #edf2fc;
}

.color {
  display: inline-block;
  vertical-align: middle;
  width: 20px;
  height: 20px;
  border: 1px solid #eaecef;
  margin-right: 4px;
}
</style>

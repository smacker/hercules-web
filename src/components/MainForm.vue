<template>
  <form class="main-form" @submit.prevent="submit">
    <input
      class="main-form__input"
      type="text"
      name="repo"
      id="repo"
      placeholder="ex: github.com/src-d/go-git"
      v-model="repo"
    >
    <input class="main-form__submit" type="submit" value="Code Burndown!" :disabled="!repo">
  </form>
</template>

<script>
const httpsRegex = new RegExp("^https?://");

export default {
  props: { onSubmit: Function },

  data() {
    return { repo: "" };
  },

  methods: {
    submit() {
      if (!this.repo) {
        return;
      }
      this.repo = this.repo.replace(httpsRegex, "");

      this.onSubmit(this.repo);
    }
  }
};
</script>


<style>
.main-form__input {
  display: block;
  width: 100%;
  margin: 0 auto 10px;
  padding: 10px;

  border: 1px solid #222;
  border-radius: 5px;

  font-size: 24px;
  text-align: center;
}

.main-form__submit {
  display: block;
  margin: 0 auto;
  padding: 10px 20px;

  border: none;
  background: #222;
  color: #fff;
  font-size: 22px;
  text-transform: uppercase;

  -webkit-font-smoothing: antialiased;

  transition: background 150ms;
  cursor: pointer;
}

.main-form__submit:hover {
  background: #90feb5;
}

.main-form__submit:disabled {
  background: #444;
  cursor: default;
}
</style>


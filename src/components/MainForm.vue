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
      <input class="main-form__submit" type="submit" value="Code Burndown!">
  </form>
</template>

<script>
export default {
  props: { onSubmit: Function },

  data() {
    return { repo: '' };
  },

  methods: {
    submit() {
      this.onSubmit(this.repo);
    }
  }
};
</script>


<style>
.main-form {
  width: 500px;
  margin: 50px auto;
}

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
</style>


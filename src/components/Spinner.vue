<template>
  <div class="spinner">loading...</div>
</template>

<style>
.spinner,
.spinner:before,
.spinner:after {
  background: #000000;
  animation: spinner-animation 1s infinite ease-in-out;
  width: 1em;
  height: 4em;
}
.spinner {
  padding-top: 2em;
  display: inline-block;
  color: #000000;
  text-indent: -9999em;
  position: relative;
  font-size: 11px;
  transform: translateZ(0);
  animation-delay: -0.16s;
}
.spinner:before,
.spinner:after {
  position: absolute;
  top: 0;
  content: "";
}
.spinner:before {
  left: -1.5em;
  animation-delay: -0.32s;
}
.spinner:after {
  left: 1.5em;
}
@-webkit-keyframes spinner-animation {
  0%,
  80%,
  100% {
    box-shadow: 0 0;
    height: 4em;
  }
  40% {
    box-shadow: 0 -2em;
    height: 5em;
  }
}
@keyframes spinner-animation {
  0%,
  80%,
  100% {
    box-shadow: 0 0;
    height: 4em;
  }
  40% {
    box-shadow: 0 -2em;
    height: 5em;
  }
}
</style>

<template>
  <span v-show="false">LoginRequired</span>
</template>

<script>
export default {
  name: 'LoginRequired',
  beforeCreate () {
    if (this.$store.state.token === '') { this.$router.push('login') }
  },
  beforeUpdate () {
    if (this.$store.state.token === '') { this.$router.push('login') }
  }
}
</script>

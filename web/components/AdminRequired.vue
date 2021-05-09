<template>
  <span v-show="false">AdminRequired</span>
</template>

<script>
export default {
  name: 'AdminRequired',
  beforeCreate () {
    if (this.$store.state.user.role === 1) { this.$router.push(-1) }
  },
  beforeUpdate () {
    if (this.$store.state.user.role === 1) { this.$router.push(-1) }
  }
}
</script>

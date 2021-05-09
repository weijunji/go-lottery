<template>
  <el-container>
    <el-header style="position: fixed; width: 100%">
      <el-row style="border-bottom: solid 1px #e6e6e6;">
        <el-col :span="18">
          <el-menu class="el-menu-demo" router mode="horizontal" style="border: none">
            <span id="logo">科软抽奖系统</span>
            <el-menu-item index="/">首页</el-menu-item>
            <el-menu-item index="/lottery">活动列表</el-menu-item>
            <el-menu-item index="/award" :disabled="!isLogin">中奖信息</el-menu-item>
            <el-menu-item index="/manage" v-if="isAdmin">抽奖管理</el-menu-item>
          </el-menu>
        </el-col>
        <el-col :span="6" style="height: 60px;">
          <span style="line-height: 60px; float: right">
            <el-dropdown v-if="isLogin" style="height: 60px;">
              <el-avatar :src="avatar" fit="fill" style="top: 50%; position: relative; transform: translateY(-50%);"></el-avatar>
              <el-dropdown-menu slot="dropdown">
                <p style="text-align: center; color: #888;">欢迎！</p>
                <el-dropdown-item>{{ username }}</el-dropdown-item>
                <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
            <router-link v-else to="/login">登录</router-link>
          </span>
        </el-col>
      </el-row>
    </el-header>
    <el-main style="padding-top: 70px;"><Nuxt /></el-main>
  </el-container>
</template>

<script>
import md5 from 'md5'

export default {
  name: 'Default',
  data () {
    return {
      clientHeight: ''
    }
  },
  computed: {
    isLogin () {
      return this.$store.state.user.token === ''
    },
    username () {
      return this.$store.state.user.username
    },
    avatar () {
      const email = this.$store.state.user.email.toLowerCase()
      return 'https://www.gravatar.com/avatar/' + md5(email)
    },
    isAdmin () {
      return this.$store.state.user.role === 0 || this.$store.state.user.role === 2
    }
  },
  method: {
    logout () {
      this.$store.commit('user/changeToken', { token: '' })
    }
  },
  mounted () {
    this.clientHeight = `${document.documentElement.clientHeight}`
    window.onresize = function () {
      this.clientHeight = `${document.documentElement.clientHeight}`
    }
  }
}
</script>

<style>
#logo {
  float: left;
  text-align: center;
  line-height: 60px;
  font-weight: bold;
  font-size: 1.2rem;
  margin-right: 10px;
}

a{
  text-decoration:none;
  color:#333;
}

html {
  font-family:
    'Source Sans Pro',
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    Roboto,
    'Helvetica Neue',
    Arial,
    sans-serif;
  font-size: 16px;
  word-spacing: 1px;
  -ms-text-size-adjust: 100%;
  -webkit-text-size-adjust: 100%;
  -moz-osx-font-smoothing: grayscale;
  -webkit-font-smoothing: antialiased;
  box-sizing: border-box;
}

*,
*::before,
*::after {
  box-sizing: border-box;
  margin: 0;
}

.button--green {
  display: inline-block;
  border-radius: 4px;
  border: 1px solid #3b8070;
  color: #3b8070;
  text-decoration: none;
  padding: 10px 30px;
}

.button--green:hover {
  color: #fff;
  background-color: #3b8070;
}

.button--grey {
  display: inline-block;
  border-radius: 4px;
  border: 1px solid #35495e;
  color: #35495e;
  text-decoration: none;
  padding: 10px 30px;
  margin-left: 15px;
}

.button--grey:hover {
  color: #fff;
  background-color: #35495e;
}
</style>

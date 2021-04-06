export const state = () => ({
  token: '',
  username: '',
  email: '',
  id: -1,
  role: -1
})

export const mutations = {
  changeToken (state, user) {
    if (user.token === '') {
      state.token = ''
      window.localStorage.token = ''
      state.username = ''
      window.localStorage.username = ''
      state.email = ''
      window.localStorage.email = ''
      state.id = -1
      window.localStorage.id = -1
      state.role = -1
      window.localStorage.role = -1
    } else {
      state.token = user.token
      window.localStorage.token = state.token
      state.username = user.username
      window.localStorage.username = state.username
      state.email = user.email
      window.localStorage.email = state.email
      state.id = user.id
      window.localStorage.id = state.id
      state.role = user.role
      window.localStorage.role = state.role
    }
  },
  init (state) {
    if (window.localStorage.token) {
      state.token = window.localStorage.token
      state.username = window.localStorage.username
      state.email = window.localStorage.email
      state.id = Number(window.localStorage.id)
      state.role = Number(window.localStorage.role)
    }
  }
}

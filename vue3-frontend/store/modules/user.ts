import { Module } from "vuex";

interface UserState {
  token: string | null;
  userInfo: {
    username: string;
    organization: string;
    isAdmin: boolean;
  } | null;
}

const userModule: Module<UserState, any> = {
  namespaced: true,
  state: {
    token: null,
    userInfo: null,
  },
  mutations: {
    SET_TOKEN(state, token: string) {
      state.token = token;
    },
    SET_USER_INFO(state, userInfo: UserState["userInfo"]) {
      state.userInfo = userInfo;
    },
    CLEAR_USER(state) {
      state.token = null;
      state.userInfo = null;
    },
  },
  actions: {
    saveUser({ commit }, { token, userInfo }) {
      commit("SET_TOKEN", token);
      commit("SET_USER_INFO", userInfo);

      // 保存到 localStorage
      localStorage.setItem("authToken", token);
      localStorage.setItem("userInfo", JSON.stringify(userInfo));
    },
    loadUser({ commit }) {
      const token = localStorage.getItem("authToken");
      const userInfo = localStorage.getItem("userInfo");

      if (token && userInfo) {
        commit("SET_TOKEN", token);
        commit("SET_USER_INFO", JSON.parse(userInfo));
      }
    },
    logout({ commit }) {
      commit("CLEAR_USER");

      // 清除 localStorage
      localStorage.removeItem("authToken");
      localStorage.removeItem("userInfo");
    },
  },
};

export default userModule;
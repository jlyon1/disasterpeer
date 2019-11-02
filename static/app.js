// this does not
new Vue({
  el: "#app",
  data() {
    return {
      message: "hello",
      uuid: "",
      messagesSaved: 0,
      name: "",
      email: "",
      phone: ""
    };
  },
  mounted() {
    fetch("/uuid")
      .then(data => {
        return data.text();
      })
      .then(val => {
        console.log(val);
        this.uuid = val;
      });
  },
  methods: {
    postInfo() {
      fetch("/info", {
        method: 'POST',
        body: JSON.stringify({
          name: this.name,
          email: this.email,
          phone: this.phone
        })
      });
    }
  }
});

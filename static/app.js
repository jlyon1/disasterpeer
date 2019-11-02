// this does not
new Vue({
  el: "#app",
  data() {
    return {
      message: "hello",
      uuid: "",
      messagesSaved: 0,
      name: "me",
      email: "me@example.com",
      phone: "555 555 5555",
      lat: -1,
      long: -1,
      locAvail: false,
      locationUpdates: 0
    };
  },
  mounted() {
    if (window.navigator !== undefined) {
      this.locAvail = true
    }
    let el = this;
    navigator.geolocation.getCurrentPosition(function(position) {
      el.lat = position.coords.latitude;
      el.long = position.coords.longitude;
      el.locationUpdates++;
      el.postInfo()
    });
    navigator.geolocation.watchPosition(function(position) {
      el.lat = position.coords.latitude;
      el.long = position.coords.longitude;
      el.locationUpdates++;
      el.postInfo()
    });

    fetch("/uuid")
      .then(data => {
        return data.text();
      })
      .then(val => {
        console.log(val);
        this.uuid = val;
      });
      this.getInfo()
      this.getAllMessages()
  },
  methods: {
    getAllMessages(){
      fetch("/messages").then((data) => data.json()).then((val) => {
        this.messagesSaved = val.length
      })
    },
    postInfo() {
      fetch("/info", {
        method: 'POST',
        body: JSON.stringify({
          name: this.name,
          email: this.email,
          phone: this.phone,
          lat: this.lat,
          long: this.long
        })
      });
    },
    getInfo() {
      fetch("/info").then((data)=> data.json()).then((val) => {
        this.name = val.Name;
        this.email = val.Email;
        this.phone = val.Phone;
      });
    }
  }
});

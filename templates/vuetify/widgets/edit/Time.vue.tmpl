<template>
  <div>
    <v-text-field :label="$attrs.label" append-icon="timer" v-model="displayTime" @click:append="openDatePicker" :placeholder="$attrs.label"></v-text-field>
    <div class="hidden-date">
      <datetime type="datetime" ref="datepicker" class="date-time-wrapper" @input="update" v-model="time" :placeholder="$attrs.label" required></datetime>
    </div>
  </div>
</template>

<script>
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import moment from 'moment';

export default {
  data() {
    return {
      time: new Date().getTime().toString()
    };
  },
  computed: {
    displayTime() {
      let time = this.time;
      return new Date(time).toLocaleString();
    }
  },
  mounted() {
    let minDate = '1970-01-01T00:00:00.000Z'
    let dateToday = moment().format("YYYY-MM-DDTHH:mm:ss.SSS");
    this.time = new Date(this.$attrs.value.toDate()).toISOString();
    if ((this.time !== null) && (this.time === minDate)) {
      this.time = dateToday;
    }
    /// converting the date from 'milliseconds' to this format : 2018-07-28T00:00:00.000Z
    console.log("setting time to : ", this.time);
    this.update();
  },
  methods: {
    update() {
      // if (!this.time.length === 0) {
      let protoDate = new google_protobuf_timestamp_pb.Timestamp();
      console.log("emiting..");
      console.log(this.time.length);
      let correct_time = new Date(this.time);

      protoDate.fromDate(new Date(this.time));
      this.$emit("gocipe", protoDate);
      console.log("sending this");
      console.log(protoDate);
      // }s
    },
    openDatePicker() {
      this.$refs["datepicker"].$el.firstChild.click();
    }
  },
  inheritAttrs: false
};
</script>

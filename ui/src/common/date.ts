const toYYYY년MM월DD일HH시MM분 = (date: Date) =>
  `${date.getFullYear()}년 ${date.getMonth() + 1}월 ${date.getDate()}일 ${date.getHours()}시 ${date.getMinutes()}분`;

export default toYYYY년MM월DD일HH시MM분;

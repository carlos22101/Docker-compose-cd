const express = require('express');
const path = require('path');
const app = express();
require('dotenv').config({ path: '../.env' });


app.use(express.static(path.join(__dirname, 'public')));
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

const PORT = 80;
app.listen(PORT, () => {
  console.log(`Frontend server running on port ${PORT}`);
});

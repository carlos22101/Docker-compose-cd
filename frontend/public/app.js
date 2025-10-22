const API_BASE = "http://34.204.46.88:8000" ;

async function fetchUsers() {
  const res = await fetch(`${API_BASE}/users`);
  return res.json();
}

async function createUser(user) {
  const res = await fetch(`${API_BASE}/users`, {
    method: 'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(user)
  });
  return res.json();
}

async function updateUser(id, user) {
  const res = await fetch(`${API_BASE}/users/${id}`, {
    method: 'PUT',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(user)
  });
  return res.json();
}

async function deleteUser(id) {
  const res = await fetch(`${API_BASE}/users/${id}`, { method:'DELETE' });
  return res.text();
}

function renderTable(users) {
  const tbody = document.querySelector('#usersTable tbody');
  tbody.innerHTML = '';
  users.forEach(u => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${u.id}</td>
      <td>${u.first_name}</td>
      <td>${u.last_name}</td>
      <td>${u.email || ''}</td>
      <td>
        <button data-id="${u.id}" class="edit">Editar</button>
        <button data-id="${u.id}" class="del">Eliminar</button>
      </td>
    `;
    tbody.appendChild(tr);
  });
}

async function loadAndRender() {
  const users = await fetchUsers();
  renderTable(users);
}

document.getElementById('userForm').addEventListener('submit', async (e) => {
  e.preventDefault();
  const id = document.getElementById('userId').value;
  const user = {
    first_name: document.getElementById('firstName').value,
    last_name: document.getElementById('lastName').value,
    email: document.getElementById('email').value
  };
  if (id) {
    await updateUser(id, user);
  } else {
    await createUser(user);
  }
  document.getElementById('userForm').reset();
  document.getElementById('userId').value = '';
  await loadAndRender();
});

document.getElementById('cancelEdit').addEventListener('click', () => {
  document.getElementById('userForm').reset();
  document.getElementById('userId').value = '';
});

document.querySelector('#usersTable tbody').addEventListener('click', async (e) => {
  if (e.target.classList.contains('edit')) {
    const id = e.target.dataset.id;
    const res = await fetch(`${API_BASE}/users/${id}`);
    const u = await res.json();
    document.getElementById('userId').value = u.id;
    document.getElementById('firstName').value = u.first_name;
    document.getElementById('lastName').value = u.last_name;
    document.getElementById('email').value = u.email || '';
  }
  if (e.target.classList.contains('del')) {
    const id = e.target.dataset.id;
    if (confirm('Eliminar usuario?')) {
      await deleteUser(id);
      await loadAndRender();
    }
  }
});


loadAndRender();

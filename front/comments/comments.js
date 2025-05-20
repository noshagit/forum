const comments = [
  { name: "David King", text: "i’m gay" },
  { name: "Feng Min", text: "May the Kitsune guide you" },
  { name: "Martin", text: "Quelqu’un a vu mon fils" },
  { name: "Don Pollo", text: "Bla bla bla ble ble ble LIIIINGAAAA GULI GULI GULIII" }
];

function createCommentElement(comment) {
  const container = document.createElement('div');
  container.className = 'comment-block';

  const avatar = document.createElement('div');
  avatar.className = 'avatar';

  const text = document.createElement('div');
  text.className = 'text';

  const name = document.createElement('div');
  name.className = 'name';
  name.textContent = comment.name;

  const desc = document.createElement('div');
  desc.className = 'desc';
  desc.textContent = comment.text;

  text.appendChild(name);
  text.appendChild(desc);

  container.appendChild(avatar);
  container.appendChild(text);

  return container;
}

function displayComments() {
  const container = document.getElementById('comments-container');
  comments.forEach(comment => {
    const el = createCommentElement(comment);
    container.appendChild(el);
  });
}

document.addEventListener('DOMContentLoaded', displayComments);

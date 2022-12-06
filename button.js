// Create a button and add it to the page
var button = document.createElement("button");
button.innerHTML = "Click me";
document.body.appendChild(button);

// Add an event listener to the button
button.addEventListener("click", function() {
  // When the button is clicked, change the title of the page
  document.title = "Button clicked!";
});

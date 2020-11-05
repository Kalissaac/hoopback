document.getElementById("user-menu").addEventListener('click', function () {
  profile.open = !profile.open
});

var openMobile = false

var mobileMenuOpen = document.getElementById('user-menu-mobile-open')
var mobileMenuClosed = document.getElementById('user-menu-mobile-closed')
var mobileMenuElement = document.getElementById('user-menu-mobile')

document.getElementById("user-menu-mobile-toggle").addEventListener('click', function () {
  if (openMobile == false) {
    mobileMenuOpen.classList.remove("hidden")
    mobileMenuOpen.classList.add("block")

    mobileMenuClosed.classList.remove("block")
    mobileMenuClosed.classList.add("hidden")

    mobileMenuElement.classList = "block md:hidden"
    openMobile = true
  } else {
    mobileMenuOpen.classList.remove("block")
    mobileMenuOpen.classList.add("hidden")

    mobileMenuClosed.classList.remove("hidden")
    mobileMenuClosed.classList.add("block")

    mobileMenuElement.className = "hidden md:hidden"
    openMobile = false
  }
});

document.getElementById("sign-out").addEventListener('click', function () {
  document.getElementById("parent-overlay").classList.remove("pointer-events-none")
  overlay.open = true
  modal.open = true
  profile.open = false
})

document.getElementById("sign-out-cancel").addEventListener('click', function () {
  document.getElementById("parent-overlay").classList.add("pointer-events-none")
  overlay.open = false
  modal.open = false
})
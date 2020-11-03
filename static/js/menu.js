var open = false

var profileElement = document.getElementById('profile')
document.getElementById("user-menu").addEventListener('click', function () {
  if (open == false) {
    profileElement.classList.remove("opacity-0")
    profileElement.classList.add("opacity-100")
    open = true
  } else {
    profileElement.classList.remove("opacity-100")
    profileElement.classList.add("opacity-0")
    open = false
  }
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
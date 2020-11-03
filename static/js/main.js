/* ===================================================================
 * Infinity - Main JS
 *
 * ------------------------------------------------------------------- */

(function($) {

	"use strict";

	var cfg = {
		defAnimation   : "fadeInUp",    // default css animation
		scrollDuration : 800,           // smoothscroll duration
	},

	$WIN = $(window);


   // Add the User Agent to the <html>
   // will be used for IE10 detection (Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0))
	var doc = document.documentElement;
	doc.setAttribute('data-useragent', navigator.userAgent);


	/* Preloader
	 * -------------------------------------------------- */
	var ssPreloader = function() {

		$WIN.on('load', function() {

			// force page scroll position to top at page refresh
			$('html, body').animate({ scrollTop: 0 }, 'normal');

	      // will first fade out the loading animation
	    	$("#loader").fadeOut("slow", function(){

	        // will fade out the whole DIV that covers the website.
	        $("#preloader").delay(300).fadeOut("slow");

	      });
	  	});
	};


  	/* OffCanvas Menu
	 * ------------------------------------------------------ */
  //  var ssOffCanvas = function() {

	//        var menuTrigger = $('#header-menu-trigger'),
	//        nav             = $('#menu-nav-wrap'),
	//        closeButton     = nav.find('.close-button'),
	//        siteBody        = $('body'),
	//        mainContents    = $('section, footer');
  //  };


  /* Smooth Scrolling
	* ------------------------------------------------------ */
	var ssSmoothScroll = function() {

		$('.smoothscroll').on('click', function (e) {
			var target = this.hash,
			$target    = $(target);

		 	e.preventDefault();
		 	e.stopPropagation();

	    	$('html, body').stop().animate({
	       	'scrollTop': $target.offset().top
	      }, cfg.scrollDuration, 'swing').promise().done(function () {

	      	// check if menu is open
	      	if ($('body').hasClass('menu-is-open')) {
					$('#header-menu-trigger').trigger('click');
				}

	      	window.location.hash = target;
	      });
	  	});

	};


  /* Placeholder Plugin Settings
	* ------------------------------------------------------ */
	var ssPlaceholder = function() {
		$('input, textarea, select').placeholder();
	};


  	/* Alert Boxes
  	------------------------------------------------------- */
  	var ssAlertBoxes = function() {

  		$('.alert-box').on('click', '.close', function() {
		  $(this).parent().fadeOut(500);
		});

  	};


  /* Animations
	* ------------------------------------------------------- */
	var ssAnimations = function() {

		if (!$("html").hasClass('no-cssanimations')) {
			$('.animate-this').waypoint({
				handler: function(direction) {

					var defAnimationEfx = cfg.defAnimation;

					if ( direction === 'down' && !$(this.element).hasClass('animated')) {
						$(this.element).addClass('item-animate');

						setTimeout(function() {
							$('body .animate-this.item-animate').each(function(ctr) {
								var el       = $(this),
								animationEfx = el.data('animate') || null;

	                  	if (!animationEfx) {
			                 	animationEfx = defAnimationEfx;
			               }

			              	setTimeout( function () {
									el.addClass(animationEfx + ' animated');
									el.removeClass('item-animate');
								}, ctr * 30);

							});
						}, 100);
					}

					// trigger once only
	       		this.destroy();
				},
				offset: '95%'
			});
		}

	};


  /* Intro Animation
	* ------------------------------------------------------- */
	var ssIntroAnimation = function() {

		$WIN.on('load', function() {

	     	if (!$("html").hasClass('no-cssanimations')) {
	     		setTimeout(function(){
	    			$('.animate-intro').each(function(ctr) {
						var el = $(this),
	                   animationEfx = el.data('animate') || null;

	               if (!animationEfx) {
	                 	animationEfx = cfg.defAnimation;
	               }

	              	setTimeout( function () {
							el.addClass(animationEfx + ' animated');
						}, ctr * 300);
					});
				}, 100);
	     	}
		});

	};


  /* Back to Top
	* ------------------------------------------------------ */
	var ssBackToTop = function() {

		var pxShow  = 500,         // height on which the button will show
		fadeInTime  = 400,         // how slow/fast you want the button to show
		fadeOutTime = 400,         // how slow/fast you want the button to hide
		scrollSpeed = 300,         // how slow/fast you want the button to scroll to top. can be a value, 'slow', 'normal' or 'fast'
		goTopButton = $("#go-top")

		// Show or hide the sticky footer button
		$(window).on('scroll', function() {
			if ($(window).scrollTop() >= pxShow) {
				goTopButton.fadeIn(fadeInTime);
			} else {
				goTopButton.fadeOut(fadeOutTime);
			}
		});
	};



  /* Initialize
	* ------------------------------------------------------ */
	(function ssInit() {

		ssPreloader();
		// ssOffCanvas();
		ssSmoothScroll();
		ssPlaceholder();
		ssAlertBoxes();
		ssAnimations();
		ssIntroAnimation();
		ssBackToTop();

	})();


})(jQuery);
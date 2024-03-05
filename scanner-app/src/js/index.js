import QRReader from './vendor/qrscan.js';
import snackbar from './snackbar.js';

import '../css/styles.css';

//If service worker is installed, show offline usage notification
if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker
        .register('/service-worker.js')
        .then((reg) => {
            // console.log('SW registered: ', reg);
            if (!localStorage.getItem('offline')) {
                localStorage.setItem('offline', true);
                snackbar.show('App is ready for offline usage.', 5000);
            }
        })
        .catch((regError) => {
            console.log('SW registration failed: ', regError);
        });
    });
}

window.addEventListener('DOMContentLoaded', () => {
    //To check the device and add iOS support
    window.iOS = ['iPad', 'iPhone', 'iPod'].indexOf(navigator.platform) >= 0;
    window.isMediaStreamAPISupported = navigator && navigator.mediaDevices && 'enumerateDevices' in navigator.mediaDevices;
    window.noCameraPermission = false;

    var frame = null;
    var scanConfirmation = document.querySelector('.app_confirmation');
    var scanConfirmationOverlay = document.querySelector('.app_confirmation-overlay');
    var scanningEle = document.querySelector('.custom-scanner');
    var appScanningEle = document.querySelector('.app__scanner-img');

    var ticketConfirmation = document.querySelector('.app_confirmation-content');
    var ticketCustomer = document.querySelector('.ticket-customer');
    var ticketEvent = document.querySelector('.ticket-event');
    var ticketCode = document.querySelector('.ticket-code');

    window.appOverlay = document.querySelector('.app__overlay');

    //Initializing qr scanner
    window.addEventListener('load', (event) => {
        // Fetch user from storage or prompt
        if (localStorage.getItem('username')) {
            window.username = localStorage.getItem('username');
        } else {
            window.username = prompt('Enter your name');
            localStorage.setItem('username', window.username);
        }

        QRReader.init(); //To initialize QR Scanner
        // Set camera overlay size
        setTimeout(() => {
            setCameraOverlay();
            if (window.isMediaStreamAPISupported) {
                scan();
            }
        }, 1000);

        // To support other browsers who dont have mediaStreamAPI
        selectFromPhoto();
    });

    function setCameraOverlay() {
        window.appOverlay.style.borderStyle = 'solid';
    }

    function createFrame() {
        frame = document.createElement('img');
        frame.src = '';
        frame.id = 'frame';
    }

    function scan(forSelectedPhotos = false) {
        if (window.isMediaStreamAPISupported && !window.noCameraPermission) {
            scanningEle.style.display = 'block';
            appScanningEle.style.display = 'block';
        }

        if (forSelectedPhotos) {
            scanningEle.style.display = 'block';
            appScanningEle.style.display = 'block';
        }

        QRReader.scan((ticketCode) => {
            validateTicket(ticketCode).then((ticket) => {
                console.log('ticket validated', ticket);
                showScanConfirmation(ticket.customer.name, ticket.eventName, ticket.valid, ticketCode);
            }).catch((err) => {
                console.log('encountered error while validating ticket', err);
                showScanConfirmation('ERROR', 'Unable to validate ticket', false, ticketCode);
            }).finally(() => {
                setTimeout(() => {
                    hideScanConfirmation();
                }, 2000);
            });
        }, forSelectedPhotos);
    }

    async function validateTicket(ticketCode) {
        var ticketParams = {
            "user": window.username,
            "ticket_code": ticketCode
        };
        const response = await fetch(
            'https://digital-venue.herokuapp.com/mobile/check-in',
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(ticketParams)
            }
        );
        return response.json().then(data => {
            // map to ticket object
            return {
                customer: data.customer,
                eventName: data.name,
                valid: data.valid
                }
        });
    }

    function showScanConfirmation(customer, eventName, valid, code) {
        scanningEle.style.display = 'none';
        appScanningEle.style.display = 'none';
        scanConfirmation.classList.remove('app_confirmation--hide');
        scanConfirmationOverlay.classList.remove('app_confirmation--hide');

        setConfirmationValues(customer, eventName, code);
        if (valid) {
            ticketConfirmation.classList.add('valid');
        } else {
            ticketConfirmation.classList.add('invalid');
        }
    }

    function setConfirmationValues(customer, eventName, code) {
        ticketCustomer.textContent = customer || '';
        ticketEvent.textContent = eventName || '';
        ticketCode.textContent = code || '';
    }

    function hideScanConfirmation() {
        if (!window.isMediaStreamAPISupported) {
            frame.src = '';
            frame.className = '';
        }

        scanConfirmation.classList.add('app_confirmation--hide');
        scanConfirmationOverlay.classList.add('app_confirmation--hide');

        setConfirmationValues();
        ticketConfirmation.classList.remove('valid');
        ticketConfirmation.classList.remove('invalid');

        scan();
    }

    function selectFromPhoto() {
        //Creating the camera element
        var camera = document.createElement('input');
        camera.setAttribute('type', 'file');
        camera.setAttribute('capture', 'camera');
        camera.id = 'camera';
        window.appOverlay.style.borderStyle = '';
        createFrame();

        //Add the camera and img element to DOM
        var pageContentElement = document.querySelector('.app__layout-content');
        pageContentElement.appendChild(camera);
        pageContentElement.appendChild(frame);

        //On camera change
        camera.addEventListener('change', (event) => {
            if (event.target && event.target.files.length > 0) {
                frame.className = 'app__overlay';
                frame.src = URL.createObjectURL(event.target.files[0]);
                if (!window.noCameraPermission) {
                    scanningEle.style.display = 'block';
                    appScanningEle.style.display = 'block';
                }
                window.appOverlay.style.borderColor = 'rgb(62, 78, 184)';
                scan(true);
            }
        });
    }
});

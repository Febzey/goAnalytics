"use strict";
// (() => {
//     "use strict";
//     /**
//      * Sending our analytic data to backend server. 
//      * data includes the event type.
//      * @param data 
//      */
//     const sendAnalyticsData = (data: any) => {
//         // Create a new Image element
//         const img = new Image();
//         // Construct the URL with data as query parameters
//         const url = `http://localhost:8080/analytics?${encodeData({
//             event: [data.event],
//             userData: [JSON.stringify(data)],
//         })}`;
//         // Set the image source to trigger the GET request
//         img.src = url;
//     };
//     // Function to encode data as query parameters
//     const encodeData = (data: any) => {
//         return Object.keys(data)
//             .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
//             .join('&');
//     };
//     const gatherUserData = () => {
//         return {
//             userAgent: navigator.userAgent,
//             url: window.location.href,
//             referrer: document.referrer,
//         };
//     };
//     const handleEvent = (eventType: string) => {
//         const userData = gatherUserData();
//         const analyticsData = {
//             event: eventType,
//             userData,
//             // Add more data as needed
//         };
//         sendAnalyticsData(analyticsData);
//     };
//     // Listen for various events
//     window.addEventListener('load', () => handleEvent('pageLoad'));
//     window.addEventListener('popstate', () => handleEvent('popstate'));
//     window.addEventListener('hashchange', () => handleEvent('hashchange'));
//     window.addEventListener('DOMContentLoaded', () => handleEvent('domcontentloaded'));
//     window.addEventListener('', () => handleEvent('domcontentloaded'));
//     document.addEventListener('click', (event) => {
//         //@ts-ignore
//         console.log("leeleek")
//     });
// })();
(() => {
    "use strict";
    const encodeData = (data) => {
        return Object.entries(data)
            .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
            .join('&');
    };
    const sendAnalyticsData = (data) => {
        const img = new Image();
        img.src = `http://localhost:8080/analytics?${encodeData({
            event: data.event,
            userData: JSON.stringify(data.userData),
        })}`;
    };
    const gatherUserData = () => {
        return {
            userAgent: navigator.userAgent,
            url: window.location.href,
            referrer: document.referrer,
        };
    };
    const handlePageLoad = () => {
        const userData = gatherUserData();
        const analyticsData = {
            event: 'pageLoad',
            userData,
        };
        sendAnalyticsData(analyticsData);
    };
    const handleRouteChange = () => {
        const userData = gatherUserData();
        const analyticsData = {
            event: 'routeChange',
            userData,
        };
        sendAnalyticsData(analyticsData);
    };
    const setupSPAListener = () => {
        const stateListener = (type) => {
            const orig = history[type];
            return function () {
                const rv = orig.apply(this, arguments);
                const event = new Event(type);
                event.arguments = arguments;
                dispatchEvent(event);
                return rv;
            };
        };
        history.pushState = stateListener('pushState');
        window.addEventListener('pushState', handleRouteChange);
        if ('onhashchange' in window) {
            window.onhashchange = handleRouteChange;
        }
    };
    // Listen for the page load event
    window.addEventListener('load', handlePageLoad);
    // Check if it's an SPA and set up appropriate listeners
    if (history.pushState !== undefined) {
        setupSPAListener();
    }
})();

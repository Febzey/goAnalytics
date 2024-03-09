(() => {
    "use strict";

    interface EventData {
        event: string
        userAgent: string;
        url: string;
        referrer: string;
        longitude?: number;
        latitude?: number;
    }

    interface AnalyticsData {
        event: string;
        EventData: EventData;
    }

    /**
     * Encoding data to send in img get.
     */
    const encodeData = (data: Record<string, string | number | boolean>): string => {
        return Object.entries(data)
            .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
            .join('&');
    };
    

    /**
     * Sending our data in an image get request to backend server.
     * @param data 
     */
    const sendAnalyticsData = (data: AnalyticsData): void => {
        const img = new Image();
        img.src = `http://localhost:8080/analytics?${encodeData({
            data: JSON.stringify(data.EventData),
        })}`;
    };

    /**
     * Getting some data on the client viewing the page.
     * @returns 
     */
    const gatherEventData = async (eventType: string): Promise<EventData> => {
        return {
            event: eventType,
            userAgent: navigator.userAgent,
            url: window.location.href,
            referrer: document.referrer,
        };
    };

    /**
     * Fires when the route is changed or page loads or hash change
     * Sends data to server.
     */
    const handleRouteChange = (type: string) => {
        gatherEventData(type)
            .then((eventData) => {
                const analyticsData: AnalyticsData = {
                    event: 'routeChange',
                    EventData: eventData,
                };
                sendAnalyticsData(analyticsData);
            })
            .catch((error) => {
                console.error('Error handling route change:', error);
            });
    };
    
    /**
     * Setting page change listener for Singe Page Applications such as
     * react, vue, angular etc.
     */
    const setupSPAListener = (): void => {
        const stateListener = (type: string) => {
            const orig = history[type as keyof History];
            return function (this: Window & typeof globalThis) {
                const rv = orig.apply(this, arguments as any);
                const event = new Event(type);
                (event as any).arguments = arguments;
                dispatchEvent(event);
                return rv;
            };
        };
        
        
        history.pushState = stateListener('pushState');
        window.addEventListener('pushState', () => handleRouteChange("pushstate"));

        if ('onhashchange' in window) {
            window.onhashchange = () => handleRouteChange("onhashchange");
        }
    };

    /**
     * Listen for when the page loads for first time.
     */
    window.addEventListener('load', () => handleRouteChange("load"));

    // Check if it's an SPA and set up appropriate listeners
    if (history.pushState !== undefined) {
        setupSPAListener();
    }
    
})();

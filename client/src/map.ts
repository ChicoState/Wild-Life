export const getLocation = () => {
    let options = {
        enableHighAccuracy: true,
        timeout: 5000,
        maximumAge: 0
    };
    let coords: any
    const success = (pos: any) => {
        coords = pos.coords;
    }
    const error = (err: any) => {
        console.warn(`ERROR(${err.code}): ${err.message}`);
    }

    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(success,error,options);
        return coords
    }
    else{
        console.log("Geolocation is not supported by this browser.");
    }
}
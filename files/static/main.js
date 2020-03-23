"use strict";

const iceConfig = {
    iceServers: [{
        urls: [
            'stun:stun1.l.google.com:19302',
            'stun:stun2.l.google.com:19302',
        ],
    }],
    iceCandidatePoolSize: 10,
};

const LaplaceVar = {
    ui: {},
};

function print(s) {
    LaplaceVar.ui.output.innerHTML += s + '\n';
}

function updateRoomUI() {
    LaplaceVar.ui.panel.style.display = 'none';
    LaplaceVar.ui.videoContainer.style.display = 'block';
    LaplaceVar.ui.streamPageUI.style.display = 'block';

    if (LaplaceVar.roomID) {
        const joinUrl = `${window.location.protocol}//${window.location.host}/?id=${LaplaceVar.roomID}`;
        LaplaceVar.ui.qrcodeObj = new QRCode(LaplaceVar.ui.qrcode, {
            text: joinUrl,
            width: 128,
            height: 128,
        });
        LaplaceVar.ui.roomText.innerHTML = '#' + LaplaceVar.roomID;
        LaplaceVar.ui.joinLinkText.innerHTML = joinUrl;
        LaplaceVar.ui.joinLinkText.href = joinUrl;
    }
}

function initUI() {
    LaplaceVar.ui.inputRoomID = document.getElementById('inputRoomID');
    LaplaceVar.ui.joinLinkText = document.getElementById("join-link");
    LaplaceVar.ui.joinForm = document.getElementById('joinForm');
    LaplaceVar.ui.output = document.getElementById('output');
    LaplaceVar.ui.qrcode = document.getElementById("qrcode");
    LaplaceVar.ui.panel = document.getElementById('panel');
    LaplaceVar.ui.roomText = document.getElementById('room-text');
    LaplaceVar.ui.streamPageUI = document.getElementById('stream-page-ui');
    LaplaceVar.ui.video = document.getElementById('mainVideo');
    LaplaceVar.ui.videoContainer = document.getElementById('video-container');

    LaplaceVar.ui.joinForm.onsubmit = (e) => {
        e.preventDefault();
        LaplaceVar.roomID = LaplaceVar.ui.inputRoomID.value;
        window.location.href = window.location.href + "?id=" + LaplaceVar.roomID;
    };
    print("Logs:");
    print("[+] Page loaded");
}

function getWebsocketUrl() {
    if (window.location.protocol === "https:") {
        return `wss://${window.location.host}`
    } else {
        return `ws://${window.location.host}`
    }
}

async function newRoom(rID) {
    print("[+] Get room ID: " + rID);
    LaplaceVar.roomID = rID;
    updateRoomUI();
}

async function newSessionStream(sessionID) {
    print('[+] New session: ' + sessionID);
    LaplaceVar.pcs[sessionID] = new RTCPeerConnection(iceConfig);
    LaplaceVar.pcs[sessionID].onicecandidate = e => {
        print('[+] Debug onicecandidate: ' + JSON.stringify(e));
        if (!e.candidate) {
            print('[+] Debug onicecandidate: got final candidate!');
            return;
        }
        print('[+] Send addCallerIceCandidate to websocket: ' + JSON.stringify(e.candidate));
        LaplaceVar.socket.send(JSON.stringify({
            Type: "addCallerIceCandidate",
            SessionID: sessionID,
            Value: JSON.stringify(e.candidate),
        }))
    };
    LaplaceVar.pcs[sessionID].oniceconnectionstatechange = () => {
        print('[+] Debug oniceconnectionstatechange ' + LaplaceVar.pcs[sessionID].iceConnectionState);
        if (LaplaceVar.pcs[sessionID].iceConnectionState === 'disconnected') {
            print("[-] Disconnected with a Peer " + sessionID);
            LaplaceVar.pcs[sessionID].close();
            LaplaceVar.pcs[sessionID] = null;
        }
    };

    LaplaceVar.mediaStream.getTracks().forEach(track => {
        LaplaceVar.pcs[sessionID].addTrack(track, LaplaceVar.mediaStream);
    });

    print('[+] Creating offer');
    const offer = await LaplaceVar.pcs[sessionID].createOffer();
    await LaplaceVar.pcs[sessionID].setLocalDescription(offer);

    print('[+] Send offer to websocket: ' + JSON.stringify(offer));
    LaplaceVar.socket.send(JSON.stringify({
        Type: "gotOffer",
        SessionID: sessionID,
        Value: JSON.stringify(offer),
    }));
}

async function addCalleeIceCandidate(sessionID, v) {
    print('[+] Debug addCalleeIceCandidate ' + sessionID + ' ' + JSON.stringify(v));
    return LaplaceVar.pcs[sessionID].addIceCandidate(v);
}

async function gotAnswer(sessionID, v) {
    print('[+] Debug gotAnswer ' + sessionID + ' ' + JSON.stringify(v));
    return LaplaceVar.pcs[sessionID].setRemoteDescription(new RTCSessionDescription(v));
}

async function doStream() {
    LaplaceVar.pcs = {}; // contains RTCPeerConnections

    updateRoomUI();

    print('[+] Initiate media: capture display media');
    // noinspection JSUnresolvedFunction
    LaplaceVar.mediaStream = await navigator.mediaDevices.getDisplayMedia({ video: true, audio: true });
    LaplaceVar.ui.video.srcObject = LaplaceVar.mediaStream;

    print('[+] Initiate websocket');
    LaplaceVar.socket = new WebSocket(getWebsocketUrl() + '/ws_serve');
    LaplaceVar.socket.onopen = async function () {
        print("[+] Connected to websocket");
    };
    LaplaceVar.socket.onmessage = async function (e) {
        print("[+] Received websocket message: " + JSON.stringify(e.data));
        try {
            const jsonData = JSON.parse(e.data);
            if (jsonData.Type === "newRoom") {
                await newRoom(jsonData.Value);
            } else if (jsonData.Type === "newSession") {
                await newSessionStream(jsonData.SessionID);
            } else if (jsonData.Type === "addCalleeIceCandidate") {
                await addCalleeIceCandidate(jsonData.SessionID, JSON.parse(jsonData.Value));
            } else if (jsonData.Type === "gotAnswer") {
                await gotAnswer(jsonData.SessionID, JSON.parse(jsonData.Value));
            }
        } catch (e) {
            print("[!] ERROR: " + e);
            console.error(e);
        }
    };
}

async function newSessionJoin(sID) {
    print('[+] New session: ' + sID);
    LaplaceVar.sessionID = sID;
    LaplaceVar.pc = new RTCPeerConnection(iceConfig);
    LaplaceVar.pc.onicecandidate = e => {
        print('[+] Debug onicecandidate: ' + JSON.stringify(e));
        if (!e.candidate) {
            print('[+] Debug onicecandidate: got final candidate!');
            return;
        }
        print('[+] Send addCalleeIceCandidate to websocket: ' + JSON.stringify(e.candidate));
        LaplaceVar.socket.send(JSON.stringify({
            Type: "addCalleeIceCandidate",
            SessionID: LaplaceVar.sessionID,
            Value: JSON.stringify(e.candidate),
        }))
    };
    LaplaceVar.pc.oniceconnectionstatechange = () => {
        print('[+] pc.oniceconnectionstatechange' + LaplaceVar.pc.iceConnectionState);
        if (LaplaceVar.pc.iceConnectionState === 'disconnected') {
            print("[-] Disconnected with Peer");
            LaplaceVar.pc.close();
            LaplaceVar.pc = null;
        }
    };
    LaplaceVar.pc.ontrack = event => {
        print('[+] Debug pc.ontrack: ' + JSON.stringify(event));
        event.streams[0].getTracks().forEach(track => {
            LaplaceVar.mediaStream.addTrack(track)
        });
        LaplaceVar.ui.video.play();
    }
}

async function addCallerIceCandidate(sID, v) {
    print('[+] Debug addCallerIceCandidate ' + sID + ' ' + JSON.stringify(v));
    if (LaplaceVar.sessionID !== sID) return;
    return LaplaceVar.pc.addIceCandidate(v);
}

async function gotOffer(sID, v) {
    print('[+] Debug gotOffer ' + sID + ' ' + JSON.stringify(v));
    if (LaplaceVar.sessionID !== sID) return;
    await LaplaceVar.pc.setRemoteDescription(new RTCSessionDescription(v));

    print('[+] Create answer');
    const answer = await LaplaceVar.pc.createAnswer();
    await LaplaceVar.pc.setLocalDescription(answer);

    print('[+] Send answer to websocket: ' + answer);
    LaplaceVar.socket.send(JSON.stringify({
        Type: "gotAnswer",
        SessionID: LaplaceVar.sessionID,
        Value: JSON.stringify(answer),
    }))
}

async function doJoin(roomID) {
    if (roomID) {
        LaplaceVar.roomID = roomID;
    } else {
        LaplaceVar.roomID = LaplaceVar.ui.inputRoomID.value;
        window.location.href = window.location.href + "?id=" + LaplaceVar.roomID;
    }
    // normalize roomID starting with #
    LaplaceVar.roomID = LaplaceVar.roomID.startsWith('#') ? LaplaceVar.roomID.slice(1) : LaplaceVar.roomID;

    updateRoomUI();

    print('[+] Initiate media: set remote source');
    LaplaceVar.mediaStream = new MediaStream();
    LaplaceVar.ui.video.srcObject = LaplaceVar.mediaStream;

    print('[+] Initiate websocket');
    LaplaceVar.socket = new WebSocket(getWebsocketUrl() + "/ws_connect?id=" + LaplaceVar.roomID);
    LaplaceVar.socket.onopen = async function () {
        print("[+] Connected to websocket");
    };
    LaplaceVar.socket.onmessage = async function (e) {
        print("[+] Received websocket message: " + JSON.stringify(e.data));
        try {
            const jsonData = JSON.parse(e.data);
            if (jsonData.Type === "newSession") {
                await newSessionJoin(jsonData.SessionID);
            } else if (jsonData.Type === "addCallerIceCandidate") {
                await addCallerIceCandidate(jsonData.SessionID, JSON.parse(jsonData.Value));
            } else if (jsonData.Type === "gotOffer") {
                await gotOffer(jsonData.SessionID, JSON.parse(jsonData.Value));
            }
        } catch (e) {
            print("[!] ERROR: " + e);
            console.error(e)
        }
    };
}

async function joinByUrl() {
    const matches = /id=([^&#=]*)/.exec(window.location.search);
    if (matches && matches.length >= 2) {
        const idParam = matches[1];
        await doJoin(idParam);
    }
}

initUI();
joinByUrl();

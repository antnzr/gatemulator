const WsEvent = {
	Message: {
		New: "message/new",
		Updated: "message/updated",
	},
	DelayedMessage: {
		New: "delayed-message/new",
		Updated: "delayed-message/updated",
		Edited: "delayed-message/edited",
		Remove: "delayed-message/remove",
	},
	Channel: {
		Updated: "channel/updated",
		TokenExpired: "channel/token-expired",
	},
	User: {
		New: "user/new",
		Deleted: "user/deleted",
		Fired: "user/fired",
		Updated: "user/updated",
		PhoneChange: "user/phone/change",
		RoleChange: "user/role/change",
	},
	Ticket: {
		New: "ticket/new",
		GroupChat: "ticket/group-chat",
		Hide: "ticket/hide",
		TicketStage: "ticket/stage",
		Responsible: "ticket/responsible",
		Updated: "ticket/updated",
		Remind: "ticket/remind",
		Taken: "ticket/taken",
		Added: "ticket/added",
		MarkAsRead: "ticket/mark-as-read",
		Deleted: "ticket/deleted",
	},
	CustomField: {
		Value: {
			Deleted: "custom-field/value/deleted",
			Updated: "custom-field/value/updated",
			Add: "custom-field/value/add",
		},
		Deleted: "custom-field/deleted",
		Updated: "custom-field/updated",
		Add: "custom-field/add",
	},
	Tags: {
		New: "tags/new",
		Deleted: "tags/deleted",
	},
	Template: {
		New: "template/new",
		Updated: "template/updated",
		Deleted: "template/deleted",
	},
	TicketStage: {
		New: "ticket-stage/new",
		Updated: "ticket-stage/updated",
		Deleted: "ticket-stage/deleted",
	},
	Payment: {
		Failed: "payment/failed",
	},
	CrmSubscription: {
		ExpireBefore: "crm-subscription/expire/before",
		Expired: "crm-subscription/expired",
		ExpireAfter: "crm-subscription/expire/after",
		Updated: "crm-subscription/updated",
	},
	TicketPreviewField: {
		Added: "ticket-preview/field/add",
		Deleted: "ticket-preview/field/deleted",
		Updated: "ticket-preview/field/updated",
	},
	Invite: {
		Sent: "invite/sent",
		Deleted: "invite/deleted",
		Updated: "invite/updated",
		Accepted: "invite/accepted",
	},
	CrmProcessUser: {
		Added: "crm-process-user/added",
		Deleted: "crm-process-user/deleted",
	},
	Contact: {
		Added: "contact/added",
		Updated: "contact/updated",
	},
	Company: {
		Updated: "company/updated",
	},
};

const functions = {};

const state = {
	ticketIds: [],
	channels: [],
};

functions.loginOnAuth = async function (context) {
	const tokensResponse = await fetch(
		`${context.vars.authServerUrl}/api/v1/auth/sign-in?x-secret-token=${context.vars.authServerToken}`,
		{
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				tzOffset: 123,
				method: "sms",
				ip: "127.0.0.1",
				provider: "crm",
				phoneNumber: context.vars.accountPhone,
				ua: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0",
			}),
		}
	);
	const tokensData = await tokensResponse.json();

	const signInResponse = await fetch(`${context.vars.target}/sign-in`, {
		method: "POST",
		body: JSON.stringify(tokensData.data),
		headers: { "Content-Type": "application/json" },
	});

	const cookie = signInResponse.headers.get("set-cookie");
	if (!cookie?.length) throw new Error("No cookies (");

	context.vars.auth_cookie = cookie;
	await _connectToSocket(context);
};

function attachListeners(socket, obj = WsEvent, path = []) {
	for (const key in obj) {
		if (typeof obj[key] === "object") {
			attachListeners(socket, obj[key], [...path, key]);
		} else {
			const eventName = obj[key];
			socket.on(eventName, (data) => {
				// console.log(JSON.stringify({ event: eventName, data }, null, 2));
			});
		}
	}
}

async function _connectToSocket(context) {
	const io = require("socket.io-client");
	const socket = io(context.vars.socketUrl, {
		forceNew: true,
		reconnection: false,
		extraHeaders: { cookie: context.vars.auth_cookie },
		transports: ["websocket", "polling"],
	});

	socket
		.on("connect_error", (err) => {
			console.log(`connect_error due to ${err.message}`);
			console.log(err.description);
			console.log(err.context);
		})
		.on("connect", async () => console.log("ws connected"))
		.on("error", (err) => console.error(err.message))
		.on("disconnect", (reason) => console.log("ws disconnect: " + reason));

	attachListeners(socket);
}

functions.extractSubscriptionId = function (
	requestParams,
	response,
	context,
	ee,
	next
) {
	const responseBody = JSON.parse(response.body);
	state.channels.push(responseBody.subscriptionId);
	return next();
};

functions.generateRandomPhone = function (context, events, done) {
	const prefix = "+770";
	const randomPart = Math.floor(100000000 + Math.random() * 900000000)
		.toString()
		.substring(0, 7);

	const phoneNumber = prefix + randomPart;
	context.vars.phoneNumber = phoneNumber;
	return done();
};

functions.getSubscriptionId = function (context, events, done) {
	context.vars.subscriptionId = state.channels.shift();
	return done();
};

functions.saveTickets = function (requestParams, response, context, ee, next) {
	const responseBody = JSON.parse(response.body);
	if (Array.isArray(responseBody.data) && responseBody.data.length) {
		state.ticketIds = responseBody.data;
	}
	return next();
};

functions.setRandomTicketId = function (context, events, done) {
	if (!state.ticketIds.length) throw new Error("No tickets");

	const randomIndex = Math.floor(Math.random() * (state.ticketIds.length - 1));

	context.vars.ticketId = state.ticketIds[randomIndex].id;
	return done();
};

module.exports = functions;

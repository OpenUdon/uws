# Browser registration call supplement 1.2

Select `x-uws-operation-profile: uws.browser-registration-call.1.2` explicitly.
The extension retains every required 1.1 field and additionally requires
`verification: reviewed_flow`. This is authority to use the selected reviewed
1.2 flow's bounded verification policy; it never authorizes another origin,
widget, trigger or application submission. The exact profile bytes, flow,
input binding, operation, approval and one-attempt claim remain bound.

Use `ValidateBrowserRegistrationCallBindingForProfile` with this discriminator
to validate the actual recipe/call pair. The unversioned call supplement and
binding APIs retain their 1.0 and 1.1 defaults respectively. A 1.1 call cannot
execute a 1.2 profile; unsupported versions fail before browser execution.
UWS core, private input envelope 1.0 and older published contracts are unchanged.

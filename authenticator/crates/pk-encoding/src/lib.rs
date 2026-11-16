//! pk-encoding: byte-to-text encodings used by the authenticator.
//!
//! Why this crate exists: the authenticator must read and write base64url and
//! hex (see plans/epics/E09-rust-crypto.md, E09-T01). This placeholder only
//! proves the Rust toolchain, workspace and test runner are wired up. The real
//! encoders arrive in E09-T01.
//!
//! `#![no_std]` means this crate does not depend on the operating system. Why:
//! the same code must later run on a microcontroller (E15). Rust's `core`
//! library is always available, even without `std`.
#![no_std]

/// Returns the crate version string.
///
/// Why a placeholder function: it gives the test something real to call, so
/// `cargo test` exercises code rather than an empty crate.
pub fn crate_name() -> &'static str {
    "pk-encoding"
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Checks the placeholder returns the expected name.
    /// `assert_eq!` panics (fails the test) if the two values differ.
    #[test]
    fn crate_name_is_pk_encoding() {
        assert_eq!(crate_name(), "pk-encoding");
    }
}

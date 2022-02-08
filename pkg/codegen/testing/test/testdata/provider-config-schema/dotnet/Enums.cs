// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.ComponentModel;
using Pulumi;

namespace Pulumi.Configstation
{
    [EnumType]
    public readonly struct Color : IEquatable<Color>
    {
        private readonly string _value;

        private Color(string value)
        {
            _value = value ?? throw new ArgumentNullException(nameof(value));
        }

        public static Color Blue { get; } = new Color("blue");
        public static Color Red { get; } = new Color("red");

        public static bool operator ==(Color left, Color right) => left.Equals(right);
        public static bool operator !=(Color left, Color right) => !left.Equals(right);

        public static explicit operator string(Color value) => value._value;

        [EditorBrowsable(EditorBrowsableState.Never)]
        public override bool Equals(object? obj) => obj is Color other && Equals(other);
        public bool Equals(Color other) => string.Equals(_value, other._value, StringComparison.Ordinal);

        [EditorBrowsable(EditorBrowsableState.Never)]
        public override int GetHashCode() => _value?.GetHashCode() ?? 0;

        public override string ToString() => _value;
    }
}
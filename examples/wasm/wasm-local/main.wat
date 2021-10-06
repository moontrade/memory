(module
  (type $0 (func (param i32)))
  (type $1 (func (param i32 i32 i32 i32) (result i32)))
  (type $2 (func (param i32 i64 i32) (result i32)))
  (type $3 (func (param i32 i32) (result i32)))
  (type $4 (func (param i64)))
  (type $5 (func))
  (type $6 (func (param i32) (result i32)))
  (type $7 (func (param i32 i64 i64)))
  (type $8 (func (param f64 i32 i32) (result f64)))
  (type $9 (func (param i32 f64 i32 i32)))
  (type $10 (func (param i32 f64)))
  (type $11 (func (param f64) (result f64)))
  (type $12 (func (param f64 f64 i32 i32) (result f64)))
  (type $13 (func (param f64 f64 i32) (result f64)))
  (type $14 (func (param f64 i32 i32 i32) (result f64)))
  (type $15 (func (param f64) (result i32)))
  (type $16 (func (param i32 i32 i32)))
  (type $17 (func (param i32 i32)))
  (type $18 (func (result i32)))
  (type $19 (func (param i32 i32 i32) (result i32)))
  (type $20 (func (param i32 i32 i32 i32 i32) (result i32)))
  (type $21 (func (result i64)))
  (import "wasi_snapshot_preview1" "fd_write" (func $28 (param i32 i32 i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "clock_time_get" (func $29 (param i32 i64 i32) (result i32)))
  (import "wasi_snapshot_preview1" "args_sizes_get" (func $30 (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "args_get" (func $31 (param i32 i32) (result i32)))
  (import "env" "setTimeout" (func $32 (param i64)))
  (memory $23  2)
  (table $22  10 10 funcref)
  (global $24  (mut i32) (i32.const 65536))
  (export "memory" (memory $23))
  (export "math.Acosh" (func $36))
  (export "math.Log" (func $37))
  (export "math.Sqrt" (func $38))
  (export "math.Log1p" (func $39))
  (export "math.Frexp" (func $40))
  (export "math.Asin" (func $42))
  (export "math.Acos" (func $45))
  (export "math.Asinh" (func $46))
  (export "math.Atan" (func $47))
  (export "math.Atan2" (func $48))
  (export "math.Atanh" (func $49))
  (export "math.Cbrt" (func $50))
  (export "math.Max" (func $51))
  (export "math.Min" (func $52))
  (export "math.Erf" (func $53))
  (export "math.Exp" (func $54))
  (export "math.Ldexp" (func $56))
  (export "math.Erfc" (func $57))
  (export "math.Exp2" (func $58))
  (export "math.Expm1" (func $59))
  (export "math.Floor" (func $60))
  (export "math.Modf" (func $61))
  (export "math.Ceil" (func $62))
  (export "math.Trunc" (func $63))
  (export "math.Pow" (func $64))
  (export "math.Sin" (func $66))
  (export "math.Hypot" (func $69))
  (export "math.Cos" (func $70))
  (export "math.Mod" (func $71))
  (export "math.Log10" (func $72))
  (export "math.Log2" (func $73))
  (export "math.Remainder" (func $74))
  (export "math.Sinh" (func $75))
  (export "math.Cosh" (func $76))
  (export "math.Tan" (func $77))
  (export "math.Tanh" (func $78))
  (export "free" (func $108))
  (export "_start" (func $117))
  (export "resume" (func $125))
  (export "malloc" (func $126))
  (export "calloc" (func $128))
  (export "realloc" (func $129))
  (export "stub" (func $137))
  (elem $25 (i32.const 1)
    $114 $115 $122 $123 $138 $139 $141 $142
    $143)
  
  (func $33 (type $5)
    )
  
  (func $34 (type $6)
    (param $0 i32)
    (result i32)
    block $block
      local.get $0
      br_if $block
      i32.const 32
      return
    end ;; $block
    i32.const 0
    local.get $0
    i32.sub
    local.get $0
    i32.and
    i32.const 125613361
    i32.mul
    i32.const 27
    i32.shr_u
    i32.const 65536
    i32.add
    i32.load8_u
    )
  
  (func $35 (type $7)
    (param $0 i32)
    (param $1 i64)
    (param $2 i64)
    (local $3 i64)
    (local $4 i64)
    local.get $0
    local.get $2
    local.get $1
    i64.mul
    i64.store offset=8
    local.get $0
    local.get $2
    i64.const 4294967295
    i64.and
    local.tee $3
    local.get $1
    i64.const 4294967295
    i64.and
    local.tee $4
    i64.mul
    i64.const 32
    i64.shr_u
    local.get $3
    local.get $1
    i64.const 32
    i64.shr_u
    local.tee $1
    i64.mul
    i64.add
    local.tee $3
    i64.const 32
    i64.shr_u
    local.get $2
    i64.const 32
    i64.shr_u
    local.tee $2
    local.get $1
    i64.mul
    i64.add
    local.get $3
    i64.const 4294967295
    i64.and
    local.get $2
    local.get $4
    i64.mul
    i64.add
    i64.const 32
    i64.shr_u
    i64.add
    i64.store
    )
  
  (func $36 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    f64.const +nan:0x8000000000001
    local.set $3
    block $block
      local.get $0
      f64.const 0x1.0000000000000p-0
      f64.lt
      br_if $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      f64.const 0x0.0000000000000p-1023
      local.set $3
      local.get $0
      f64.const 0x1.0000000000000p-0
      f64.eq
      br_if $block
      block $block_0
        local.get $0
        f64.const 0x1.0000000000000p+28
        f64.ge
        i32.const 1
        i32.xor
        br_if $block_0
        local.get $0
        local.get $4
        local.get $4
        call $37
        f64.const 0x1.62e42fefa39efp-1
        f64.add
        return
      end ;; $block_0
      block $block_1
        local.get $0
        f64.const 0x1.0000000000000p+1
        f64.gt
        i32.const 1
        i32.xor
        br_if $block_1
        local.get $0
        local.get $0
        f64.add
        f64.const -0x1.0000000000000p-0
        local.get $0
        local.get $0
        f64.mul
        f64.const -0x1.0000000000000p-0
        f64.add
        local.get $4
        local.get $4
        call $38
        local.get $0
        f64.add
        f64.div
        f64.add
        local.get $4
        local.get $4
        call $37
        return
      end ;; $block_1
      local.get $0
      f64.const -0x1.0000000000000p-0
      f64.add
      local.set $0
      local.get $0
      local.get $0
      local.get $0
      f64.add
      local.get $0
      local.get $0
      f64.mul
      f64.add
      local.get $4
      local.get $4
      call $38
      f64.add
      local.get $4
      local.get $4
      call $39
      local.set $3
    end ;; $block
    local.get $3
    )
  
  (func $37 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    (local $4 f64)
    (local $5 i32)
    (local $6 f64)
    (local $7 f64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $3
    global.set $24
    block $block
      block $block_0
        local.get $0
        local.get $0
        f64.ne
        br_if $block_0
        local.get $0
        f64.const 0x1.fffffffffffffp+1023
        f64.gt
        br_if $block_0
        f64.const +nan:0x8000000000001
        local.set $4
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        br_if $block
        f64.const -inf
        local.set $4
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.eq
        br_if $block
        local.get $3
        local.get $0
        local.get $3
        local.get $3
        call $40
        local.get $3
        i32.load offset=8
        local.get $3
        f64.load
        local.tee $0
        f64.const 0x1.6a09e667f3bcdp-1
        f64.lt
        local.tee $5
        i32.sub
        f64.convert_i32_s
        local.tee $6
        f64.const 0x1.62e42fee00000p-1
        f64.mul
        local.get $0
        local.get $0
        f64.add
        local.get $0
        local.get $5
        select
        f64.const -0x1.0000000000000p-0
        f64.add
        local.tee $4
        local.get $4
        f64.const 0x1.0000000000000p-1
        f64.mul
        f64.mul
        local.tee $7
        local.get $6
        f64.const 0x1.a39ef35793c76p-33
        f64.mul
        local.get $4
        local.get $4
        f64.const 0x1.0000000000000p+1
        f64.add
        f64.div
        local.tee $0
        local.get $7
        local.get $0
        local.get $0
        f64.mul
        local.tee $6
        local.get $6
        f64.mul
        local.tee $0
        local.get $0
        local.get $0
        f64.const 0x1.39a09d078c69fp-3
        f64.mul
        f64.const 0x1.c71c51d8e78afp-3
        f64.add
        f64.mul
        f64.const 0x1.999999997fa04p-2
        f64.add
        f64.mul
        local.get $6
        local.get $0
        local.get $0
        local.get $0
        f64.const 0x1.2f112df3e5244p-3
        f64.mul
        f64.const 0x1.7466496cb03dep-3
        f64.add
        f64.mul
        f64.const 0x1.2492494229359p-2
        f64.add
        f64.mul
        f64.const 0x1.5555555555593p-1
        f64.add
        f64.mul
        f64.add
        f64.add
        f64.mul
        f64.add
        f64.sub
        local.get $4
        f64.sub
        f64.sub
        local.set $4
        br $block
      end ;; $block_0
      local.get $0
      local.set $4
    end ;; $block
    local.get $3
    i32.const 16
    i32.add
    global.set $24
    local.get $4
    )
  
  (func $38 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i64)
    (local $5 i32)
    (local $6 i32)
    (local $7 i64)
    (local $8 i64)
    (local $9 i64)
    (local $10 i64)
    block $block
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.eq
      local.get $0
      local.get $0
      f64.ne
      i32.or
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      f64.const +nan:0x8000000000001
      local.set $3
      block $block_0
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        br_if $block_0
        block $block_1
          local.get $0
          i64.reinterpret_f64
          local.tee $4
          i64.const 52
          i64.shr_u
          i32.wrap_i64
          i32.const 2047
          i32.and
          local.tee $5
          br_if $block_1
          i32.const 1
          local.set $5
          loop $loop
            local.get $4
            i64.const 4503599627370496
            i64.and
            i64.const 0
            i64.ne
            br_if $block_1
            local.get $5
            i32.const -1
            i32.add
            local.set $5
            local.get $4
            i64.const 1
            i64.shl
            local.set $4
            br $loop
          end ;; $loop
        end ;; $block_1
        local.get $4
        i64.const -9218868437227405313
        i64.and
        i64.const 4503599627370496
        i64.or
        local.get $5
        i32.const -1023
        i32.add
        local.tee $6
        i32.const 1
        i32.and
        i64.extend_i32_u
        i64.shl
        local.set $7
        i64.const 9007199254740992
        local.set $4
        i64.const 0
        local.set $8
        i64.const 0
        local.set $9
        block $block_2
          loop $loop_0
            local.get $7
            i64.const 1
            i64.shl
            local.set $7
            local.get $4
            i64.eqz
            br_if $block_2
            local.get $9
            local.get $4
            local.get $9
            i64.add
            local.tee $10
            local.get $4
            i64.add
            local.get $7
            local.get $10
            i64.lt_u
            local.tee $5
            select
            local.set $9
            local.get $7
            i64.const 0
            local.get $10
            local.get $5
            select
            i64.sub
            local.set $7
            i64.const 0
            local.get $4
            local.get $5
            select
            local.get $8
            i64.add
            local.set $8
            local.get $4
            i64.const 1
            i64.shr_u
            local.set $4
            br $loop_0
          end ;; $loop_0
        end ;; $block_2
        i64.const 0
        local.get $8
        i64.const 1
        i64.and
        local.get $7
        i64.eqz
        select
        local.get $8
        i64.add
        i64.const 1
        i64.shr_u
        local.get $6
        i32.const 1
        i32.shr_u
        i32.const 1022
        i32.add
        i64.extend_i32_u
        i64.const 52
        i64.shl
        i64.add
        f64.reinterpret_i64
        local.set $3
      end ;; $block_0
      local.get $3
      return
    end ;; $block
    local.get $0
    )
  
  (func $39 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i64)
    (local $5 f64)
    (local $6 i32)
    (local $7 f64)
    (local $8 i64)
    (local $9 f64)
    (local $10 f64)
    f64.const +nan:0x8000000000001
    local.set $3
    block $block
      local.get $0
      f64.const -0x1.0000000000000p-0
      f64.lt
      br_if $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      f64.const -inf
      local.set $3
      local.get $0
      f64.const -0x1.0000000000000p-0
      f64.eq
      br_if $block
      f64.const +inf
      local.set $3
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      block $block_0
        block $block_1
          local.get $0
          i64.reinterpret_f64
          local.tee $4
          i64.const 9223372036854775807
          i64.and
          f64.reinterpret_i64
          local.tee $3
          f64.const 0x1.a827999fcef32p-2
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_1
          block $block_2
            local.get $3
            f64.const 0x1.0000000000000p-29
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_2
            block $block_3
              local.get $3
              f64.const 0x1.0000000000000p-54
              f64.lt
              i32.eqz
              br_if $block_3
              local.get $0
              return
            end ;; $block_3
            local.get $0
            local.get $0
            local.get $0
            f64.mul
            f64.const -0x1.0000000000000p-1
            f64.mul
            f64.add
            return
          end ;; $block_2
          local.get $0
          f64.const -0x1.2bec333018867p-2
          f64.gt
          i32.const 1
          i32.xor
          br_if $block_1
          local.get $0
          f64.const 0x1.0000000000000p-1
          f64.mul
          local.get $0
          f64.mul
          local.set $5
          i32.const 0
          local.set $6
          f64.const 0x0.0000000000000p-1023
          local.set $7
          br $block_0
        end ;; $block_1
        block $block_4
          block $block_5
            local.get $3
            f64.const 0x1.0000000000000p+53
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_5
            local.get $0
            local.get $0
            f64.const 0x1.0000000000000p-0
            f64.add
            local.tee $3
            f64.sub
            f64.const 0x1.0000000000000p-0
            f64.add
            local.get $0
            local.get $3
            f64.const -0x1.0000000000000p-0
            f64.add
            f64.sub
            local.get $3
            i64.reinterpret_f64
            local.tee $4
            i64.const 52
            i64.shr_u
            i32.wrap_i64
            local.tee $6
            i32.const 1023
            i32.gt_u
            select
            local.get $3
            f64.div
            local.set $7
            br $block_4
          end ;; $block_5
          local.get $4
          i64.const 52
          i64.shr_u
          i32.wrap_i64
          local.set $6
          f64.const 0x0.0000000000000p-1023
          local.set $7
        end ;; $block_4
        block $block_6
          block $block_7
            local.get $4
            i64.const 4503599627370495
            i64.and
            local.tee $4
            i64.const 1865452045155276
            i64.gt_u
            br_if $block_7
            local.get $4
            i64.const 4607182418800017408
            i64.or
            local.set $8
            local.get $6
            i32.const -1023
            i32.add
            local.set $6
            br $block_6
          end ;; $block_7
          local.get $4
          i64.const 4602678819172646912
          i64.or
          local.set $8
          local.get $6
          i32.const -1022
          i32.add
          local.set $6
          i64.const 4503599627370496
          local.get $4
          i64.sub
          i64.const 2
          i64.shr_u
          local.set $4
        end ;; $block_6
        local.get $8
        f64.reinterpret_i64
        f64.const -0x1.0000000000000p-0
        f64.add
        local.tee $0
        local.get $0
        f64.const 0x1.0000000000000p-1
        f64.mul
        f64.mul
        local.set $5
        local.get $4
        i64.const 0
        i64.ne
        br_if $block_0
        f64.const 0x0.0000000000000p-1023
        local.set $3
        block $block_8
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.ne
          br_if $block_8
          local.get $6
          i32.eqz
          br_if $block
          local.get $6
          f64.convert_i32_s
          local.tee $0
          f64.const 0x1.62e42fee00000p-1
          f64.mul
          local.get $7
          local.get $0
          f64.const 0x1.a39ef35793c76p-33
          f64.mul
          f64.add
          f64.add
          return
        end ;; $block_8
        local.get $5
        local.get $0
        f64.const -0x1.5555555555555p-1
        f64.mul
        f64.const 0x1.0000000000000p-0
        f64.add
        f64.mul
        local.set $3
        block $block_9
          local.get $6
          br_if $block_9
          local.get $0
          local.get $3
          f64.sub
          return
        end ;; $block_9
        local.get $6
        f64.convert_i32_s
        local.tee $5
        f64.const 0x1.62e42fee00000p-1
        f64.mul
        local.get $3
        local.get $7
        local.get $5
        f64.const 0x1.a39ef35793c76p-33
        f64.mul
        f64.add
        f64.sub
        local.get $0
        f64.sub
        f64.sub
        return
      end ;; $block_0
      local.get $0
      local.get $0
      f64.const 0x1.0000000000000p+1
      f64.add
      f64.div
      local.tee $9
      local.get $9
      f64.mul
      local.tee $3
      local.get $3
      local.get $3
      local.get $3
      local.get $3
      local.get $3
      local.get $3
      f64.const 0x1.2f112df3e5244p-3
      f64.mul
      f64.const 0x1.39a09d078c69fp-3
      f64.add
      f64.mul
      f64.const 0x1.7466496cb03dep-3
      f64.add
      f64.mul
      f64.const 0x1.c71c51d8e78afp-3
      f64.add
      f64.mul
      f64.const 0x1.2492494229359p-2
      f64.add
      f64.mul
      f64.const 0x1.999999997fa04p-2
      f64.add
      f64.mul
      f64.const 0x1.5555555555593p-1
      f64.add
      f64.mul
      local.set $3
      block $block_10
        local.get $6
        br_if $block_10
        local.get $0
        local.get $5
        local.get $9
        local.get $5
        local.get $3
        f64.add
        f64.mul
        f64.sub
        f64.sub
        return
      end ;; $block_10
      local.get $6
      f64.convert_i32_s
      local.tee $10
      f64.const 0x1.62e42fee00000p-1
      f64.mul
      local.get $5
      local.get $7
      local.get $10
      f64.const 0x1.a39ef35793c76p-33
      f64.mul
      f64.add
      local.get $9
      local.get $5
      local.get $3
      f64.add
      f64.mul
      f64.add
      f64.sub
      local.get $0
      f64.sub
      f64.sub
      local.set $3
    end ;; $block
    local.get $3
    )
  
  (func $40 (type $9)
    (param $0 i32)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $4
    global.set $24
    i32.const 0
    local.set $5
    block $block
      local.get $1
      f64.const 0x0.0000000000000p-1023
      f64.eq
      br_if $block
      local.get $1
      local.get $1
      f64.ne
      br_if $block
      local.get $1
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $1
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $4
      local.get $1
      call $41
      local.get $4
      i32.load offset=8
      local.get $4
      i64.load
      local.tee $6
      i64.const 52
      i64.shr_u
      i32.wrap_i64
      i32.const 2047
      i32.and
      i32.add
      i32.const -1022
      i32.add
      local.set $5
      local.get $6
      i64.const -9218868437227405313
      i64.and
      i64.const 4602678819172646912
      i64.or
      f64.reinterpret_i64
      local.set $1
    end ;; $block
    local.get $0
    local.get $5
    i32.store offset=8
    local.get $0
    local.get $1
    f64.store
    local.get $4
    i32.const 16
    i32.add
    global.set $24
    )
  
  (func $41 (type $10)
    (param $0 i32)
    (param $1 f64)
    block $block
      local.get $1
      i64.reinterpret_f64
      i64.const 9223372036854775807
      i64.and
      f64.reinterpret_i64
      f64.const 0x1.0000000000000p-1022
      f64.lt
      i32.const 1
      i32.xor
      br_if $block
      local.get $0
      i32.const -52
      i32.store offset=8
      local.get $0
      local.get $1
      f64.const 0x1.0000000000000p+52
      f64.mul
      f64.store
      return
    end ;; $block
    local.get $0
    i32.const 0
    i32.store offset=8
    local.get $0
    local.get $1
    f64.store
    )
  
  (func $42 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 f64)
    (local $5 i32)
    block $block
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.ne
      br_if $block
      local.get $0
      return
    end ;; $block
    f64.const +nan:0x8000000000001
    local.set $3
    block $block_0
      local.get $0
      f64.neg
      local.get $0
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      select
      local.tee $4
      f64.const 0x1.0000000000000p-0
      f64.gt
      br_if $block_0
      f64.const 0x1.0000000000000p-0
      local.get $4
      local.get $4
      f64.mul
      f64.sub
      local.get $5
      local.get $5
      call $38
      local.set $3
      block $block_1
        block $block_2
          local.get $4
          f64.const 0x1.6666666666666p-1
          f64.gt
          i32.const 1
          i32.xor
          br_if $block_2
          f64.const 0x1.921fb54442d18p-0
          local.get $3
          local.get $4
          f64.div
          call $43
          f64.sub
          local.set $4
          br $block_1
        end ;; $block_2
        local.get $4
        local.get $3
        f64.div
        call $43
        local.set $4
      end ;; $block_1
      local.get $4
      f64.neg
      local.get $4
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      select
      local.set $3
    end ;; $block_0
    local.get $3
    )
  
  (func $43 (type $11)
    (param $0 f64)
    (result f64)
    block $block
      block $block_0
        local.get $0
        f64.const 0x1.51eb851eb851fp-1
        f64.le
        i32.const 1
        i32.xor
        i32.eqz
        br_if $block_0
        local.get $0
        f64.const 0x1.3504f333f9de6p+1
        f64.gt
        i32.const 1
        i32.xor
        i32.eqz
        br_if $block
        local.get $0
        f64.const -0x1.0000000000000p-0
        f64.add
        local.get $0
        f64.const 0x1.0000000000000p-0
        f64.add
        f64.div
        call $44
        f64.const 0x1.921fb54442d18p-1
        f64.add
        f64.const 0x1.1a62633145c07p-55
        f64.add
        return
      end ;; $block_0
      local.get $0
      call $44
      return
    end ;; $block
    f64.const 0x1.921fb54442d18p-0
    f64.const 0x1.0000000000000p-0
    local.get $0
    f64.div
    call $44
    f64.sub
    f64.const 0x1.1a62633145c07p-54
    f64.add
    )
  
  (func $44 (type $11)
    (param $0 f64)
    (result f64)
    (local $1 f64)
    local.get $0
    local.get $0
    f64.mul
    local.tee $1
    local.get $1
    local.get $1
    local.get $1
    local.get $1
    f64.const -0x1.c007fa1f72594p-1
    f64.mul
    f64.const -0x1.028545b6b807ap+4
    f64.add
    f64.mul
    f64.const -0x1.2c08c36880273p+6
    f64.add
    f64.mul
    f64.const -0x1.eb8bf2d05ba25p+6
    f64.add
    f64.mul
    f64.const -0x1.03669fd28ec8ep+6
    f64.add
    f64.mul
    local.get $1
    local.get $1
    local.get $1
    local.get $1
    local.get $1
    f64.const 0x1.8dbc45b14603cp+4
    f64.add
    f64.mul
    f64.const 0x1.4a0dd43b8fa25p+7
    f64.add
    f64.mul
    f64.const 0x1.b0e18d2e2be3bp+8
    f64.add
    f64.mul
    f64.const 0x1.e563f13b049eap+8
    f64.add
    f64.mul
    f64.const 0x1.8519efbbd62ecp+7
    f64.add
    f64.div
    local.get $0
    f64.mul
    local.get $0
    f64.add
    )
  
  (func $45 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    f64.const 0x1.921fb54442d18p-0
    local.get $0
    local.get $3
    local.get $3
    call $42
    f64.sub
    )
  
  (func $46 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    (local $5 f64)
    block $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      block $block_0
        block $block_1
          local.get $0
          f64.neg
          local.get $0
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          select
          local.tee $3
          f64.const 0x1.0000000000000p+28
          f64.gt
          i32.const 1
          i32.xor
          br_if $block_1
          local.get $3
          local.get $4
          local.get $4
          call $37
          f64.const 0x1.62e42fefa39efp-1
          f64.add
          local.set $3
          br $block_0
        end ;; $block_1
        block $block_2
          local.get $3
          f64.const 0x1.0000000000000p+1
          f64.gt
          i32.const 1
          i32.xor
          br_if $block_2
          local.get $3
          local.get $3
          f64.add
          f64.const 0x1.0000000000000p-0
          local.get $3
          local.get $3
          local.get $3
          f64.mul
          f64.const 0x1.0000000000000p-0
          f64.add
          local.get $4
          local.get $4
          call $38
          f64.add
          f64.div
          f64.add
          local.get $4
          local.get $4
          call $37
          local.set $3
          br $block_0
        end ;; $block_2
        local.get $3
        f64.const 0x1.0000000000000p-28
        f64.lt
        br_if $block_0
        local.get $3
        local.get $3
        f64.mul
        local.set $5
        local.get $3
        local.get $5
        local.get $5
        f64.const 0x1.0000000000000p-0
        f64.add
        local.get $4
        local.get $4
        call $38
        f64.const 0x1.0000000000000p-0
        f64.add
        f64.div
        f64.add
        local.get $4
        local.get $4
        call $39
        local.set $3
      end ;; $block_0
      local.get $3
      f64.neg
      local.get $3
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      select
      local.set $0
    end ;; $block
    local.get $0
    )
  
  (func $47 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    block $block
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.eq
      br_if $block
      block $block_0
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.gt
        i32.const 1
        i32.xor
        br_if $block_0
        local.get $0
        call $43
        return
      end ;; $block_0
      local.get $0
      f64.neg
      call $43
      f64.neg
      local.set $0
    end ;; $block
    local.get $0
    )
  
  (func $48 (type $12)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 i64)
    (local $5 i32)
    block $block
      local.get $0
      local.get $0
      f64.ne
      local.get $1
      local.get $1
      f64.ne
      i32.or
      i32.eqz
      br_if $block
      f64.const +nan:0x8000000000001
      return
    end ;; $block
    block $block_0
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.ne
      br_if $block_0
      local.get $0
      i64.reinterpret_f64
      i64.const -9223372036854775808
      i64.and
      local.set $4
      block $block_1
        local.get $1
        f64.const 0x0.0000000000000p-1023
        f64.ge
        i32.const 1
        i32.xor
        br_if $block_1
        local.get $1
        i64.reinterpret_f64
        i64.const 0
        i64.lt_s
        br_if $block_1
        local.get $4
        f64.reinterpret_i64
        return
      end ;; $block_1
      local.get $4
      i64.const 4614256656552045848
      i64.or
      f64.reinterpret_i64
      return
    end ;; $block_0
    block $block_2
      local.get $1
      f64.const 0x0.0000000000000p-1023
      f64.ne
      br_if $block_2
      local.get $0
      i64.reinterpret_f64
      i64.const -9223372036854775808
      i64.and
      i64.const 4609753056924675352
      i64.or
      f64.reinterpret_i64
      return
    end ;; $block_2
    block $block_3
      block $block_4
        block $block_5
          local.get $1
          f64.const 0x1.fffffffffffffp+1023
          f64.gt
          br_if $block_5
          local.get $0
          f64.const 0x1.fffffffffffffp+1023
          f64.gt
          local.get $0
          f64.const -0x1.fffffffffffffp+1023
          f64.lt
          i32.or
          local.set $5
          local.get $1
          f64.const -0x1.fffffffffffffp+1023
          f64.lt
          i32.eqz
          br_if $block_3
          local.get $0
          i64.reinterpret_f64
          i64.const -9223372036854775808
          i64.and
          local.set $4
          local.get $5
          i32.eqz
          br_if $block_4
          local.get $4
          i64.const 4612488097114038738
          i64.or
          f64.reinterpret_i64
          return
        end ;; $block_5
        local.get $0
        i64.reinterpret_f64
        i64.const -9223372036854775808
        i64.and
        local.set $4
        block $block_6
          block $block_7
            local.get $0
            f64.const 0x1.fffffffffffffp+1023
            f64.gt
            br_if $block_7
            local.get $0
            f64.const -0x1.fffffffffffffp+1023
            f64.lt
            i32.eqz
            br_if $block_6
          end ;; $block_7
          local.get $4
          i64.const 4605249457297304856
          i64.or
          f64.reinterpret_i64
          return
        end ;; $block_6
        local.get $4
        f64.reinterpret_i64
        return
      end ;; $block_4
      local.get $4
      i64.const 4614256656552045848
      i64.or
      f64.reinterpret_i64
      return
    end ;; $block_3
    block $block_8
      local.get $5
      i32.eqz
      br_if $block_8
      local.get $0
      i64.reinterpret_f64
      i64.const -9223372036854775808
      i64.and
      i64.const 4609753056924675352
      i64.or
      f64.reinterpret_i64
      return
    end ;; $block_8
    local.get $0
    local.get $1
    f64.div
    local.get $5
    local.get $5
    call $47
    local.set $0
    block $block_9
      local.get $1
      f64.const 0x0.0000000000000p-1023
      f64.lt
      i32.const 1
      i32.xor
      br_if $block_9
      block $block_10
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.le
        i32.const 1
        i32.xor
        br_if $block_10
        local.get $0
        f64.const 0x1.921fb54442d18p+1
        f64.add
        return
      end ;; $block_10
      local.get $0
      f64.const -0x1.921fb54442d18p+1
      f64.add
      local.set $0
    end ;; $block_9
    local.get $0
    )
  
  (func $49 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    f64.const +nan:0x8000000000001
    local.set $3
    block $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      local.get $0
      f64.const -0x1.0000000000000p-0
      f64.lt
      br_if $block
      local.get $0
      f64.const 0x1.0000000000000p-0
      f64.gt
      br_if $block
      f64.const +inf
      local.set $3
      local.get $0
      f64.const 0x1.0000000000000p-0
      f64.eq
      br_if $block
      f64.const -inf
      local.set $3
      local.get $0
      f64.const -0x1.0000000000000p-0
      f64.eq
      br_if $block
      block $block_0
        local.get $0
        f64.neg
        local.get $0
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        local.tee $4
        select
        local.tee $0
        f64.const 0x1.0000000000000p-28
        f64.lt
        br_if $block_0
        local.get $0
        local.get $0
        f64.add
        local.set $3
        block $block_1
          block $block_2
            local.get $0
            f64.const 0x1.0000000000000p-1
            f64.lt
            i32.const 1
            i32.xor
            i32.eqz
            br_if $block_2
            local.get $3
            f64.const 0x1.0000000000000p-0
            local.get $0
            f64.sub
            f64.div
            local.set $0
            br $block_1
          end ;; $block_2
          local.get $3
          local.get $0
          local.get $3
          f64.mul
          f64.const 0x1.0000000000000p-0
          local.get $0
          f64.sub
          f64.div
          f64.add
          local.set $0
        end ;; $block_1
        local.get $0
        local.get $4
        local.get $4
        call $39
        f64.const 0x1.0000000000000p-1
        f64.mul
        local.set $0
      end ;; $block_0
      local.get $0
      f64.neg
      local.get $0
      local.get $4
      select
      local.set $3
    end ;; $block
    local.get $3
    )
  
  (func $50 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    (local $4 f64)
    (local $5 i64)
    (local $6 f64)
    block $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.eq
      local.get $0
      local.get $0
      f64.ne
      i32.or
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $0
      f64.neg
      local.get $0
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      local.tee $3
      select
      local.tee $4
      i64.reinterpret_f64
      i64.const 3
      i64.div_u
      local.set $5
      block $block_0
        block $block_1
          local.get $4
          f64.const 0x1.0000000000000p-1022
          f64.lt
          i32.const 1
          i32.xor
          i32.eqz
          br_if $block_1
          local.get $5
          i64.const 3071306043645493248
          i64.add
          local.set $5
          br $block_0
        end ;; $block_1
        local.get $4
        f64.const 0x1.0000000000000p+54
        f64.mul
        i64.reinterpret_f64
        i64.const 3
        i64.div_u
        i64.const 2990241250352824320
        i64.add
        local.set $5
      end ;; $block_0
      local.get $4
      f64.const 0x1.9b6db6db6db6ep-0
      local.get $5
      f64.reinterpret_i64
      local.tee $0
      local.get $0
      f64.mul
      local.get $4
      f64.div
      local.get $0
      f64.mul
      f64.const 0x1.15f15f15f15f1p-1
      f64.add
      local.tee $6
      f64.const 0x1.6a0ea0ea0ea0fp-0
      f64.add
      f64.const -0x1.691de2532c834p-1
      local.get $6
      f64.div
      f64.add
      f64.div
      f64.const 0x1.6db6db6db6db7p-2
      f64.add
      local.get $0
      f64.mul
      i64.reinterpret_f64
      i64.const 1073741824
      i64.add
      i64.const -1073741824
      i64.and
      f64.reinterpret_i64
      local.tee $0
      local.get $0
      f64.mul
      f64.div
      local.tee $4
      local.get $0
      f64.sub
      local.get $0
      local.get $0
      f64.add
      local.get $4
      f64.add
      f64.div
      local.get $0
      f64.mul
      local.get $0
      f64.add
      local.tee $0
      f64.neg
      local.get $0
      local.get $3
      select
      local.set $0
    end ;; $block
    local.get $0
    )
  
  (func $51 (type $12)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 f64)
    f64.const +inf
    local.set $4
    block $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $1
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      block $block_0
        local.get $0
        local.get $0
        f64.ne
        local.get $1
        local.get $1
        f64.ne
        i32.or
        i32.eqz
        br_if $block_0
        f64.const +nan:0x8000000000001
        return
      end ;; $block_0
      block $block_1
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.ne
        br_if $block_1
        local.get $0
        local.get $1
        f64.ne
        br_if $block_1
        local.get $1
        local.get $0
        local.get $0
        i64.reinterpret_f64
        i64.const 0
        i64.lt_s
        select
        return
      end ;; $block_1
      local.get $0
      local.get $1
      local.get $0
      local.get $1
      f64.gt
      select
      local.set $4
    end ;; $block
    local.get $4
    )
  
  (func $52 (type $12)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 f64)
    f64.const -inf
    local.set $4
    block $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $1
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      block $block_0
        local.get $0
        local.get $0
        f64.ne
        local.get $1
        local.get $1
        f64.ne
        i32.or
        i32.eqz
        br_if $block_0
        f64.const +nan:0x8000000000001
        return
      end ;; $block_0
      block $block_1
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.ne
        br_if $block_1
        local.get $0
        local.get $1
        f64.ne
        br_if $block_1
        local.get $0
        local.get $1
        local.get $0
        i64.reinterpret_f64
        i64.const 0
        i64.lt_s
        select
        return
      end ;; $block_1
      local.get $0
      local.get $1
      local.get $0
      local.get $1
      f64.lt
      select
      local.set $4
    end ;; $block
    local.get $4
    )
  
  (func $53 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 f64)
    (local $5 f64)
    (local $6 f64)
    (local $7 f64)
    (local $8 f64)
    (local $9 f64)
    (local $10 f64)
    (local $11 f64)
    (local $12 f64)
    (local $13 i32)
    block $block
      local.get $0
      local.get $0
      f64.eq
      br_if $block
      f64.const +nan:0x8000000000001
      return
    end ;; $block
    f64.const 0x1.0000000000000p-0
    local.set $3
    block $block_0
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block_0
      f64.const -0x1.0000000000000p-0
      local.set $3
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block_0
      block $block_1
        block $block_2
          local.get $0
          f64.neg
          local.get $0
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          select
          local.tee $3
          f64.const 0x1.b000000000000p-1
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_2
          block $block_3
            local.get $3
            f64.const 0x1.0000000000000p-28
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_3
            block $block_4
              local.get $3
              f64.const 0x1.0000000000000p-1015
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_4
              local.get $3
              f64.const 0x1.0000000000000p+3
              f64.mul
              local.get $3
              f64.const 0x1.06eba8214db69p-0
              f64.mul
              f64.add
              f64.const 0x1.0000000000000p-3
              f64.mul
              local.set $3
              br $block_1
            end ;; $block_4
            local.get $3
            local.get $3
            f64.const 0x1.06eba8214db69p-3
            f64.mul
            f64.add
            local.set $3
            br $block_1
          end ;; $block_3
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          f64.mul
          local.tee $4
          local.get $4
          local.get $4
          local.get $4
          f64.const -0x1.8ead6120016acp-16
          f64.mul
          f64.const -0x1.7a291236668e4p-8
          f64.add
          f64.mul
          f64.const -0x1.d2a51dbd7194fp-6
          f64.add
          f64.mul
          f64.const -0x1.4cd7d691cb913p-2
          f64.add
          f64.mul
          f64.const 0x1.06eba8214db68p-3
          f64.add
          local.get $4
          local.get $4
          local.get $4
          local.get $4
          local.get $4
          f64.const -0x1.09c4342a26120p-18
          f64.mul
          f64.const 0x1.15dc9221c1a10p-13
          f64.add
          f64.mul
          f64.const 0x1.4d022c4d36b0fp-8
          f64.add
          f64.mul
          f64.const 0x1.0a54c5536cebap-4
          f64.add
          f64.mul
          f64.const 0x1.97779cddadc09p-2
          f64.add
          f64.mul
          f64.const 0x1.0000000000000p-0
          f64.add
          f64.div
          f64.mul
          f64.add
          local.set $3
          br $block_1
        end ;; $block_2
        block $block_5
          block $block_6
            block $block_7
              local.get $3
              f64.const 0x1.4000000000000p-0
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_7
              local.get $3
              f64.const -0x1.0000000000000p-0
              f64.add
              local.tee $3
              local.get $3
              local.get $3
              local.get $3
              local.get $3
              local.get $3
              f64.const -0x1.1bf380a96073fp-9
              f64.mul
              f64.const 0x1.22a36599795ebp-5
              f64.add
              f64.mul
              f64.const -0x1.c63983d3e28ecp-4
              f64.add
              f64.mul
              f64.const 0x1.45fca805120e4p-2
              f64.add
              f64.mul
              f64.const -0x1.7d240fbb8c3f1p-2
              f64.add
              f64.mul
              f64.const 0x1.a8d00ad92b34dp-2
              f64.add
              f64.mul
              f64.const -0x1.359b8bef77538p-9
              f64.add
              local.get $3
              local.get $3
              local.get $3
              local.get $3
              local.get $3
              local.get $3
              f64.const 0x1.88b545735151dp-7
              f64.mul
              f64.const 0x1.bedc26b51dd1cp-7
              f64.add
              f64.mul
              f64.const 0x1.02660e763351fp-3
              f64.add
              f64.mul
              f64.const 0x1.2635cd99fe9a7p-4
              f64.add
              f64.mul
              f64.const 0x1.14af092eb6f33p-1
              f64.add
              f64.mul
              f64.const 0x1.b3e6618eee323p-4
              f64.add
              f64.mul
              f64.const 0x1.0000000000000p-0
              f64.add
              f64.div
              local.set $3
              local.get $0
              f64.const 0x0.0000000000000p-1023
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_6
              f64.const -0x1.b0ac160000000p-1
              local.get $3
              f64.sub
              return
            end ;; $block_7
            block $block_8
              local.get $3
              f64.const 0x1.8000000000000p+2
              f64.ge
              i32.const 1
              i32.xor
              br_if $block_8
              f64.const -0x1.0000000000000p-0
              f64.const 0x1.0000000000000p-0
              local.get $0
              f64.const 0x0.0000000000000p-1023
              f64.lt
              select
              return
            end ;; $block_8
            f64.const 0x1.0000000000000p-0
            local.get $3
            local.get $3
            f64.mul
            f64.div
            local.set $4
            block $block_9
              local.get $3
              f64.const 0x1.6db6db6db6db7p+1
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_9
              local.get $4
              local.get $4
              f64.const -0x1.eeff2ee749a62p-5
              f64.mul
              f64.const 0x1.a47ef8e484a93p+2
              f64.add
              f64.mul
              f64.const 0x1.b28a3ee48ae2cp+6
              f64.add
              local.set $5
              local.get $4
              local.get $4
              local.get $4
              local.get $4
              local.get $4
              local.get $4
              local.get $4
              f64.const -0x1.3a0efc69ac25cp+3
              f64.mul
              f64.const -0x1.4526557e4d2f2p+6
              f64.add
              f64.mul
              f64.const -0x1.7135cebccabb2p+7
              f64.add
              f64.mul
              f64.const -0x1.44cb184282266p+7
              f64.add
              f64.mul
              f64.const -0x1.f300ae4cba38dp+5
              f64.add
              f64.mul
              f64.const -0x1.51e0441b0e726p+3
              f64.add
              f64.mul
              f64.const -0x1.63416e4ba7360p-1
              f64.add
              f64.mul
              f64.const -0x1.43412600d6435p-7
              f64.add
              local.set $6
              f64.const 0x1.3a6b9bd707687p+4
              local.set $7
              f64.const 0x1.1350c526ae721p+7
              local.set $8
              f64.const 0x1.b290dd58a1a71p+8
              local.set $9
              f64.const 0x1.42b1921ec2868p+9
              local.set $10
              f64.const 0x1.ad02157700314p+8
              local.set $11
              br $block_5
            end ;; $block_9
            local.get $4
            f64.const -0x1.670e242712d62p+4
            f64.mul
            f64.const 0x1.da874e79fe763p+8
            f64.add
            local.set $5
            local.get $4
            local.get $4
            local.get $4
            local.get $4
            local.get $4
            local.get $4
            f64.const -0x1.e384e9bdc383fp+8
            f64.mul
            f64.const -0x1.004616a2e5992p+10
            f64.add
            f64.mul
            f64.const -0x1.3ec881375f228p+9
            f64.add
            f64.mul
            f64.const -0x1.4145d43c5ed98p+7
            f64.add
            f64.mul
            f64.const -0x1.1c209555f995ap+4
            f64.add
            f64.mul
            f64.const -0x1.993ba70c285dep-1
            f64.add
            f64.mul
            f64.const -0x1.4341239e86f4ap-7
            f64.add
            local.set $6
            f64.const 0x1.e568b261d5190p+4
            local.set $7
            f64.const 0x1.45cae221b9f0ap+8
            local.set $8
            f64.const 0x1.802eb189d5118p+10
            local.set $9
            f64.const 0x1.8ffb7688c246ap+11
            local.set $10
            f64.const 0x1.3f219cedf3be6p+11
            local.set $11
            br $block_5
          end ;; $block_6
          local.get $3
          f64.const 0x1.b0ac160000000p-1
          f64.add
          local.set $3
          br $block_0
        end ;; $block_5
        f64.const -0x1.2000000000000p-1
        local.get $3
        i64.reinterpret_f64
        i64.const -4294967296
        i64.and
        f64.reinterpret_i64
        local.tee $12
        local.get $12
        f64.mul
        f64.sub
        local.get $13
        local.get $13
        call $54
        local.get $12
        local.get $3
        f64.sub
        local.get $3
        local.get $12
        f64.add
        f64.mul
        local.get $6
        local.get $4
        local.get $4
        local.get $4
        local.get $4
        local.get $4
        local.get $4
        local.get $5
        f64.mul
        local.get $11
        f64.add
        f64.mul
        local.get $10
        f64.add
        f64.mul
        local.get $9
        f64.add
        f64.mul
        local.get $8
        f64.add
        f64.mul
        local.get $7
        f64.add
        f64.mul
        f64.const 0x1.0000000000000p-0
        f64.add
        f64.div
        f64.add
        local.get $13
        local.get $13
        call $54
        f64.mul
        local.get $3
        f64.div
        local.set $3
        block $block_10
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_10
          local.get $3
          f64.const -0x1.0000000000000p-0
          f64.add
          return
        end ;; $block_10
        f64.const 0x1.0000000000000p-0
        local.get $3
        f64.sub
        return
      end ;; $block_1
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      i32.const 1
      i32.xor
      br_if $block_0
      local.get $3
      f64.neg
      return
    end ;; $block_0
    local.get $3
    )
  
  (func $54 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    block $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      f64.const 0x0.0000000000000p-1023
      local.set $3
      block $block_0
        local.get $0
        f64.const -0x1.fffffffffffffp+1023
        f64.lt
        br_if $block_0
        f64.const +inf
        local.set $3
        local.get $0
        f64.const 0x1.62e42fefa39efp+9
        f64.gt
        br_if $block_0
        f64.const 0x0.0000000000000p-1023
        local.set $3
        local.get $0
        f64.const -0x1.74910d52d3051p+9
        f64.lt
        br_if $block_0
        block $block_1
          local.get $0
          f64.const -0x1.0000000000000p-28
          f64.gt
          i32.const 1
          i32.xor
          br_if $block_1
          local.get $0
          f64.const 0x1.0000000000000p-28
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_1
          local.get $0
          f64.const 0x1.0000000000000p-0
          f64.add
          return
        end ;; $block_1
        block $block_2
          block $block_3
            local.get $0
            f64.const 0x0.0000000000000p-1023
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_3
            i32.const 0
            i32.const 2147483647
            i32.const -2147483648
            local.get $0
            f64.const 0x1.71547652b82fep-0
            f64.mul
            f64.const -0x1.0000000000000p-1
            f64.add
            local.tee $3
            f64.const -0x1.0000000000000p+31
            f64.ge
            local.tee $4
            select
            local.get $3
            local.get $3
            f64.ne
            select
            local.set $5
            local.get $3
            f64.const 0x1.fffffffc00000p+30
            f64.le
            local.set $6
            block $block_4
              block $block_5
                local.get $3
                f64.abs
                f64.const 0x1.0000000000000p+31
                f64.lt
                i32.eqz
                br_if $block_5
                local.get $3
                i32.trunc_f64_s
                local.set $7
                br $block_4
              end ;; $block_5
              i32.const -2147483648
              local.set $7
            end ;; $block_4
            local.get $7
            local.get $5
            local.get $6
            select
            local.get $5
            local.get $4
            select
            local.set $5
            br $block_2
          end ;; $block_3
          i32.const 0
          local.set $5
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.gt
          i32.const 1
          i32.xor
          br_if $block_2
          i32.const 0
          i32.const 2147483647
          i32.const -2147483648
          local.get $0
          f64.const 0x1.71547652b82fep-0
          f64.mul
          f64.const 0x1.0000000000000p-1
          f64.add
          local.tee $3
          f64.const -0x1.0000000000000p+31
          f64.ge
          local.tee $4
          select
          local.get $3
          local.get $3
          f64.ne
          select
          local.set $5
          local.get $3
          f64.const 0x1.fffffffc00000p+30
          f64.le
          local.set $6
          block $block_6
            block $block_7
              local.get $3
              f64.abs
              f64.const 0x1.0000000000000p+31
              f64.lt
              i32.eqz
              br_if $block_7
              local.get $3
              i32.trunc_f64_s
              local.set $7
              br $block_6
            end ;; $block_7
            i32.const -2147483648
            local.set $7
          end ;; $block_6
          local.get $7
          local.get $5
          local.get $6
          select
          local.get $5
          local.get $4
          select
          local.set $5
        end ;; $block_2
        local.get $0
        local.get $5
        f64.convert_i32_s
        local.tee $3
        f64.const -0x1.62e42fee00000p-1
        f64.mul
        f64.add
        local.get $3
        f64.const 0x1.a39ef35793c76p-33
        f64.mul
        local.get $5
        call $55
        local.set $3
      end ;; $block_0
      local.get $3
      return
    end ;; $block
    local.get $0
    )
  
  (func $55 (type $13)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 f64)
    local.get $0
    local.get $1
    local.get $0
    local.get $1
    f64.sub
    local.tee $3
    local.get $3
    local.get $3
    local.get $3
    f64.mul
    local.tee $4
    local.get $4
    local.get $4
    local.get $4
    local.get $4
    f64.const 0x1.6376972bea4d0p-25
    f64.mul
    f64.const -0x1.bbd41c5d26bf1p-20
    f64.add
    f64.mul
    f64.const 0x1.1566aaf25de2cp-14
    f64.add
    f64.mul
    f64.const -0x1.6c16c16bebd93p-9
    f64.add
    f64.mul
    f64.const 0x1.5555555555555p-3
    f64.add
    f64.mul
    f64.sub
    local.tee $4
    f64.mul
    f64.const 0x1.0000000000000p+1
    local.get $4
    f64.sub
    f64.div
    f64.sub
    f64.sub
    f64.const 0x1.0000000000000p-0
    f64.add
    local.get $2
    local.get $2
    local.get $2
    call $56
    )
  
  (func $56 (type $14)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 i32)
    (local $5 i64)
    (local $6 i32)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $4
    global.set $24
    block $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.eq
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $4
      local.get $0
      call $41
      block $block_0
        local.get $1
        local.get $4
        i32.load offset=8
        i32.add
        local.get $4
        f64.load
        local.tee $0
        i64.reinterpret_f64
        local.tee $5
        i64.const 52
        i64.shr_u
        i32.wrap_i64
        i32.const 2047
        i32.and
        i32.add
        i32.const -1023
        i32.add
        local.tee $1
        i32.const -1076
        i32.gt_s
        br_if $block_0
        local.get $5
        i64.const -9223372036854775808
        i64.and
        f64.reinterpret_i64
        local.set $0
        br $block
      end ;; $block_0
      block $block_1
        local.get $1
        i32.const 1024
        i32.lt_s
        br_if $block_1
        f64.const -inf
        f64.const +inf
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        select
        local.set $0
        br $block
      end ;; $block_1
      f64.const 0x1.0000000000000p-53
      f64.const 0x1.0000000000000p-0
      local.get $1
      i32.const -1022
      i32.lt_s
      local.tee $6
      select
      local.get $1
      i32.const 53
      i32.add
      local.get $1
      local.get $6
      select
      i32.const 1023
      i32.add
      i64.extend_i32_u
      i64.const 52
      i64.shl
      local.get $5
      i64.const -9218868437227405313
      i64.and
      i64.or
      f64.reinterpret_i64
      f64.mul
      local.set $0
    end ;; $block
    local.get $4
    i32.const 16
    i32.add
    global.set $24
    local.get $0
    )
  
  (func $57 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 f64)
    (local $5 f64)
    (local $6 f64)
    (local $7 f64)
    (local $8 f64)
    (local $9 f64)
    (local $10 f64)
    (local $11 f64)
    (local $12 f64)
    (local $13 i32)
    block $block
      local.get $0
      local.get $0
      f64.eq
      br_if $block
      f64.const +nan:0x8000000000001
      return
    end ;; $block
    f64.const 0x0.0000000000000p-1023
    local.set $3
    block $block_0
      block $block_1
        local.get $0
        f64.const 0x1.fffffffffffffp+1023
        f64.gt
        br_if $block_1
        f64.const 0x1.0000000000000p+1
        local.set $3
        local.get $0
        f64.const -0x1.fffffffffffffp+1023
        f64.lt
        br_if $block_1
        block $block_2
          local.get $0
          f64.neg
          local.get $0
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          select
          local.tee $4
          f64.const 0x1.b000000000000p-1
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_2
          local.get $4
          f64.const 0x1.0000000000000p-56
          f64.lt
          i32.const 1
          i32.xor
          i32.eqz
          br_if $block_0
          local.get $4
          local.get $4
          local.get $4
          f64.mul
          local.tee $3
          local.get $3
          local.get $3
          local.get $3
          f64.const -0x1.8ead6120016acp-16
          f64.mul
          f64.const -0x1.7a291236668e4p-8
          f64.add
          f64.mul
          f64.const -0x1.d2a51dbd7194fp-6
          f64.add
          f64.mul
          f64.const -0x1.4cd7d691cb913p-2
          f64.add
          f64.mul
          f64.const 0x1.06eba8214db68p-3
          f64.add
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          f64.const -0x1.09c4342a26120p-18
          f64.mul
          f64.const 0x1.15dc9221c1a10p-13
          f64.add
          f64.mul
          f64.const 0x1.4d022c4d36b0fp-8
          f64.add
          f64.mul
          f64.const 0x1.0a54c5536cebap-4
          f64.add
          f64.mul
          f64.const 0x1.97779cddadc09p-2
          f64.add
          f64.mul
          f64.const 0x1.0000000000000p-0
          f64.add
          f64.div
          f64.mul
          local.set $3
          block $block_3
            local.get $4
            f64.const 0x1.0000000000000p-2
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_3
            local.get $4
            local.get $3
            f64.add
            local.set $4
            br $block_0
          end ;; $block_3
          local.get $4
          f64.const -0x1.0000000000000p-1
          f64.add
          local.get $3
          f64.add
          f64.const 0x1.0000000000000p-1
          f64.add
          local.set $4
          br $block_0
        end ;; $block_2
        block $block_4
          local.get $4
          f64.const 0x1.4000000000000p-0
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_4
          local.get $4
          f64.const -0x1.0000000000000p-0
          f64.add
          local.tee $3
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          f64.const -0x1.1bf380a96073fp-9
          f64.mul
          f64.const 0x1.22a36599795ebp-5
          f64.add
          f64.mul
          f64.const -0x1.c63983d3e28ecp-4
          f64.add
          f64.mul
          f64.const 0x1.45fca805120e4p-2
          f64.add
          f64.mul
          f64.const -0x1.7d240fbb8c3f1p-2
          f64.add
          f64.mul
          f64.const 0x1.a8d00ad92b34dp-2
          f64.add
          f64.mul
          f64.const -0x1.359b8bef77538p-9
          f64.add
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          local.get $3
          f64.const 0x1.88b545735151dp-7
          f64.mul
          f64.const 0x1.bedc26b51dd1cp-7
          f64.add
          f64.mul
          f64.const 0x1.02660e763351fp-3
          f64.add
          f64.mul
          f64.const 0x1.2635cd99fe9a7p-4
          f64.add
          f64.mul
          f64.const 0x1.14af092eb6f33p-1
          f64.add
          f64.mul
          f64.const 0x1.b3e6618eee323p-4
          f64.add
          f64.mul
          f64.const 0x1.0000000000000p-0
          f64.add
          f64.div
          local.set $3
          block $block_5
            local.get $0
            f64.const 0x0.0000000000000p-1023
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_5
            local.get $3
            f64.const 0x1.d8560b0000000p-0
            f64.add
            return
          end ;; $block_5
          f64.const 0x1.3d4fa80000000p-3
          local.get $3
          f64.sub
          return
        end ;; $block_4
        block $block_6
          block $block_7
            local.get $4
            f64.const 0x1.c000000000000p+4
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_7
            f64.const 0x1.0000000000000p-0
            local.get $4
            local.get $4
            f64.mul
            f64.div
            local.set $5
            block $block_8
              local.get $4
              f64.const 0x1.6db6db6db6db7p+1
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_8
              local.get $5
              local.get $5
              f64.const -0x1.eeff2ee749a62p-5
              f64.mul
              f64.const 0x1.a47ef8e484a93p+2
              f64.add
              f64.mul
              f64.const 0x1.b28a3ee48ae2cp+6
              f64.add
              local.set $6
              local.get $5
              local.get $5
              local.get $5
              local.get $5
              local.get $5
              local.get $5
              local.get $5
              f64.const -0x1.3a0efc69ac25cp+3
              f64.mul
              f64.const -0x1.4526557e4d2f2p+6
              f64.add
              f64.mul
              f64.const -0x1.7135cebccabb2p+7
              f64.add
              f64.mul
              f64.const -0x1.44cb184282266p+7
              f64.add
              f64.mul
              f64.const -0x1.f300ae4cba38dp+5
              f64.add
              f64.mul
              f64.const -0x1.51e0441b0e726p+3
              f64.add
              f64.mul
              f64.const -0x1.63416e4ba7360p-1
              f64.add
              f64.mul
              f64.const -0x1.43412600d6435p-7
              f64.add
              local.set $7
              f64.const 0x1.3a6b9bd707687p+4
              local.set $8
              f64.const 0x1.1350c526ae721p+7
              local.set $9
              f64.const 0x1.b290dd58a1a71p+8
              local.set $10
              f64.const 0x1.42b1921ec2868p+9
              local.set $11
              f64.const 0x1.ad02157700314p+8
              local.set $12
              br $block_6
            end ;; $block_8
            block $block_9
              local.get $0
              f64.const 0x0.0000000000000p-1023
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_9
              local.get $4
              f64.const 0x1.8000000000000p+2
              f64.gt
              br_if $block_1
            end ;; $block_9
            local.get $5
            f64.const -0x1.670e242712d62p+4
            f64.mul
            f64.const 0x1.da874e79fe763p+8
            f64.add
            local.set $6
            local.get $5
            local.get $5
            local.get $5
            local.get $5
            local.get $5
            local.get $5
            f64.const -0x1.e384e9bdc383fp+8
            f64.mul
            f64.const -0x1.004616a2e5992p+10
            f64.add
            f64.mul
            f64.const -0x1.3ec881375f228p+9
            f64.add
            f64.mul
            f64.const -0x1.4145d43c5ed98p+7
            f64.add
            f64.mul
            f64.const -0x1.1c209555f995ap+4
            f64.add
            f64.mul
            f64.const -0x1.993ba70c285dep-1
            f64.add
            f64.mul
            f64.const -0x1.4341239e86f4ap-7
            f64.add
            local.set $7
            f64.const 0x1.e568b261d5190p+4
            local.set $8
            f64.const 0x1.45cae221b9f0ap+8
            local.set $9
            f64.const 0x1.802eb189d5118p+10
            local.set $10
            f64.const 0x1.8ffb7688c246ap+11
            local.set $11
            f64.const 0x1.3f219cedf3be6p+11
            local.set $12
            br $block_6
          end ;; $block_7
          f64.const 0x1.0000000000000p+1
          f64.const 0x0.0000000000000p-1023
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          select
          local.set $3
          br $block_1
        end ;; $block_6
        f64.const -0x1.2000000000000p-1
        local.get $4
        i64.reinterpret_f64
        i64.const -4294967296
        i64.and
        f64.reinterpret_i64
        local.tee $3
        local.get $3
        f64.mul
        f64.sub
        local.get $13
        local.get $13
        call $54
        local.get $3
        local.get $4
        f64.sub
        local.get $4
        local.get $3
        f64.add
        f64.mul
        local.get $7
        local.get $5
        local.get $5
        local.get $5
        local.get $5
        local.get $5
        local.get $5
        local.get $6
        f64.mul
        local.get $12
        f64.add
        f64.mul
        local.get $11
        f64.add
        f64.mul
        local.get $10
        f64.add
        f64.mul
        local.get $9
        f64.add
        f64.mul
        local.get $8
        f64.add
        f64.mul
        f64.const 0x1.0000000000000p-0
        f64.add
        f64.div
        f64.add
        local.get $13
        local.get $13
        call $54
        f64.mul
        local.get $4
        f64.div
        local.set $3
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        i32.const 1
        i32.xor
        br_if $block_1
        f64.const 0x1.0000000000000p+1
        local.get $3
        f64.sub
        return
      end ;; $block_1
      local.get $3
      return
    end ;; $block_0
    block $block_10
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      i32.const 1
      i32.xor
      br_if $block_10
      local.get $4
      f64.const 0x1.0000000000000p-0
      f64.add
      return
    end ;; $block_10
    f64.const 0x1.0000000000000p-0
    local.get $4
    f64.sub
    )
  
  (func $58 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    block $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      f64.const 0x0.0000000000000p-1023
      local.set $3
      block $block_0
        local.get $0
        f64.const -0x1.fffffffffffffp+1023
        f64.lt
        br_if $block_0
        f64.const +inf
        local.set $3
        local.get $0
        f64.const 0x1.fffffffffffffp+9
        f64.gt
        br_if $block_0
        f64.const 0x0.0000000000000p-1023
        local.set $3
        local.get $0
        f64.const -0x1.0c80000000000p+10
        f64.lt
        br_if $block_0
        block $block_1
          block $block_2
            local.get $0
            f64.const 0x0.0000000000000p-1023
            f64.gt
            i32.const 1
            i32.xor
            br_if $block_2
            i32.const 0
            i32.const 2147483647
            i32.const -2147483648
            local.get $0
            f64.const 0x1.0000000000000p-1
            f64.add
            local.tee $3
            f64.const -0x1.0000000000000p+31
            f64.ge
            local.tee $4
            select
            local.get $3
            local.get $3
            f64.ne
            select
            local.set $5
            local.get $3
            f64.const 0x1.fffffffc00000p+30
            f64.le
            local.set $6
            block $block_3
              block $block_4
                local.get $3
                f64.abs
                f64.const 0x1.0000000000000p+31
                f64.lt
                i32.eqz
                br_if $block_4
                local.get $3
                i32.trunc_f64_s
                local.set $7
                br $block_3
              end ;; $block_4
              i32.const -2147483648
              local.set $7
            end ;; $block_3
            local.get $7
            local.get $5
            local.get $6
            select
            local.get $5
            local.get $4
            select
            local.set $5
            br $block_1
          end ;; $block_2
          i32.const 0
          local.set $5
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_1
          i32.const 0
          i32.const 2147483647
          i32.const -2147483648
          local.get $0
          f64.const -0x1.0000000000000p-1
          f64.add
          local.tee $3
          f64.const -0x1.0000000000000p+31
          f64.ge
          local.tee $4
          select
          local.get $3
          local.get $3
          f64.ne
          select
          local.set $5
          local.get $3
          f64.const 0x1.fffffffc00000p+30
          f64.le
          local.set $6
          block $block_5
            block $block_6
              local.get $3
              f64.abs
              f64.const 0x1.0000000000000p+31
              f64.lt
              i32.eqz
              br_if $block_6
              local.get $3
              i32.trunc_f64_s
              local.set $7
              br $block_5
            end ;; $block_6
            i32.const -2147483648
            local.set $7
          end ;; $block_5
          local.get $7
          local.get $5
          local.get $6
          select
          local.get $5
          local.get $4
          select
          local.set $5
        end ;; $block_1
        local.get $0
        local.get $5
        f64.convert_i32_s
        f64.sub
        local.tee $0
        f64.const 0x1.62e42fee00000p-1
        f64.mul
        local.get $0
        f64.const -0x1.a39ef35793c76p-33
        f64.mul
        local.get $5
        call $55
        local.set $3
      end ;; $block_0
      local.get $3
      return
    end ;; $block
    local.get $0
    )
  
  (func $59 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    (local $5 f64)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    (local $9 f64)
    (local $10 f64)
    (local $11 i64)
    block $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      f64.const -0x1.0000000000000p-0
      local.set $3
      block $block_0
        local.get $0
        f64.const -0x1.fffffffffffffp+1023
        f64.lt
        br_if $block_0
        block $block_1
          local.get $0
          f64.neg
          local.get $0
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          local.tee $4
          select
          local.tee $5
          f64.const 0x1.3687a9f1af2b1p+5
          f64.ge
          i32.const 1
          i32.xor
          br_if $block_1
          local.get $4
          br_if $block_0
          f64.const +inf
          local.set $3
          local.get $5
          f64.const 0x1.62e42fefa39efp+9
          f64.ge
          br_if $block_0
        end ;; $block_1
        block $block_2
          block $block_3
            block $block_4
              local.get $5
              f64.const 0x1.62e42fefa39efp-2
              f64.gt
              i32.const 1
              i32.xor
              br_if $block_4
              block $block_5
                local.get $5
                f64.const 0x1.0a2b23f3bab73p-0
                f64.lt
                i32.const 1
                i32.xor
                br_if $block_5
                block $block_6
                  local.get $0
                  f64.const 0x0.0000000000000p-1023
                  f64.lt
                  i32.const 1
                  i32.xor
                  br_if $block_6
                  local.get $0
                  f64.const 0x1.62e42fee00000p-1
                  f64.add
                  local.set $3
                  f64.const -0x1.a39ef35793c76p-33
                  local.set $5
                  i32.const -1
                  local.set $4
                  br $block_3
                end ;; $block_6
                local.get $0
                f64.const -0x1.62e42fee00000p-1
                f64.add
                local.set $3
                f64.const 0x1.a39ef35793c76p-33
                local.set $5
                i32.const 1
                local.set $4
                br $block_3
              end ;; $block_5
              i32.const 0
              i32.const 2147483647
              i32.const -2147483648
              local.get $0
              f64.const 0x1.71547652b82fep-0
              f64.mul
              local.tee $3
              f64.const -0x1.0000000000000p-1
              f64.add
              local.get $3
              f64.const 0x1.0000000000000p-1
              f64.add
              local.get $0
              f64.const 0x0.0000000000000p-1023
              f64.lt
              select
              local.tee $3
              f64.const -0x1.0000000000000p+31
              f64.ge
              local.tee $6
              select
              local.get $3
              local.get $3
              f64.ne
              select
              local.set $4
              local.get $3
              f64.const 0x1.fffffffc00000p+30
              f64.le
              local.set $7
              block $block_7
                block $block_8
                  local.get $3
                  f64.abs
                  f64.const 0x1.0000000000000p+31
                  f64.lt
                  i32.eqz
                  br_if $block_8
                  local.get $3
                  i32.trunc_f64_s
                  local.set $8
                  br $block_7
                end ;; $block_8
                i32.const -2147483648
                local.set $8
              end ;; $block_7
              local.get $8
              local.get $4
              local.get $7
              select
              local.get $4
              local.get $6
              select
              local.tee $4
              f64.convert_i32_s
              local.tee $3
              f64.const 0x1.a39ef35793c76p-33
              f64.mul
              local.set $5
              local.get $0
              local.get $3
              f64.const -0x1.62e42fee00000p-1
              f64.mul
              f64.add
              local.set $3
              br $block_3
            end ;; $block_4
            local.get $5
            f64.const 0x1.0000000000000p-54
            f64.lt
            i32.const 1
            i32.xor
            i32.eqz
            br_if $block
            i32.const 0
            local.set $4
            f64.const 0x0.0000000000000p-1023
            local.set $5
            br $block_2
          end ;; $block_3
          local.get $3
          local.get $3
          local.get $5
          f64.sub
          local.tee $0
          f64.sub
          local.get $5
          f64.sub
          local.set $5
        end ;; $block_2
        local.get $0
        local.get $0
        f64.const 0x1.0000000000000p-1
        f64.mul
        local.tee $9
        f64.mul
        local.tee $3
        local.get $3
        local.get $3
        local.get $3
        local.get $3
        local.get $3
        f64.const -0x1.afdb76e09c32dp-23
        f64.mul
        f64.const 0x1.0cfca86e65239p-18
        f64.add
        f64.mul
        f64.const -0x1.4ce199eaadbb7p-14
        f64.add
        f64.mul
        f64.const 0x1.a01a019fe5585p-10
        f64.add
        f64.mul
        f64.const -0x1.11111111110f4p-5
        f64.add
        f64.mul
        f64.const 0x1.0000000000000p-0
        f64.add
        local.tee $10
        f64.const 0x1.8000000000000p+1
        local.get $9
        local.get $10
        f64.mul
        f64.sub
        local.tee $9
        f64.sub
        f64.const 0x1.8000000000000p+2
        local.get $0
        local.get $9
        f64.mul
        f64.sub
        f64.div
        f64.mul
        local.set $9
        block $block_9
          local.get $4
          br_if $block_9
          local.get $0
          local.get $0
          local.get $9
          f64.mul
          local.get $3
          f64.sub
          f64.sub
          return
        end ;; $block_9
        local.get $0
        local.get $9
        local.get $5
        f64.sub
        f64.mul
        local.get $5
        f64.sub
        local.get $3
        f64.sub
        local.set $3
        block $block_10
          block $block_11
            block $block_12
              local.get $4
              i32.const 1
              i32.add
              br_table
                $block_10 $block_11 $block_12
                $block_11 ;; default
            end ;; $block_12
            block $block_13
              local.get $0
              f64.const -0x1.0000000000000p-2
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_13
              local.get $3
              local.get $0
              f64.const 0x1.0000000000000p-1
              f64.add
              f64.sub
              f64.const -0x1.0000000000000p+1
              f64.mul
              return
            end ;; $block_13
            local.get $0
            local.get $3
            f64.sub
            local.tee $0
            local.get $0
            f64.add
            f64.const 0x1.0000000000000p-0
            f64.add
            return
          end ;; $block_11
          block $block_14
            local.get $4
            i32.const 57
            i32.lt_u
            br_if $block_14
            local.get $4
            i64.extend_i32_u
            i64.const 52
            i64.shl
            local.get $0
            local.get $3
            f64.sub
            f64.const 0x1.0000000000000p-0
            f64.add
            i64.reinterpret_f64
            i64.add
            f64.reinterpret_i64
            f64.const -0x1.0000000000000p-0
            f64.add
            return
          end ;; $block_14
          block $block_15
            local.get $4
            i32.const 19
            i32.gt_s
            br_if $block_15
            local.get $4
            i64.extend_i32_u
            local.tee $11
            i64.const 52
            i64.shl
            i64.const 4607182418800017408
            i64.const 9007199254740992
            local.get $11
            i64.shr_u
            i64.sub
            f64.reinterpret_i64
            local.get $3
            local.get $0
            f64.sub
            f64.sub
            i64.reinterpret_f64
            i64.add
            f64.reinterpret_i64
            return
          end ;; $block_15
          local.get $4
          i64.extend_i32_u
          i64.const 52
          i64.shl
          local.get $0
          local.get $3
          i32.const 1023
          local.get $4
          i32.sub
          i64.extend_i32_u
          i64.const 52
          i64.shl
          f64.reinterpret_i64
          f64.add
          f64.sub
          f64.const 0x1.0000000000000p-0
          f64.add
          i64.reinterpret_f64
          i64.add
          f64.reinterpret_i64
          return
        end ;; $block_10
        local.get $0
        local.get $3
        f64.sub
        f64.const 0x1.0000000000000p-1
        f64.mul
        f64.const -0x1.0000000000000p-1
        f64.add
        local.set $3
      end ;; $block_0
      local.get $3
      return
    end ;; $block
    local.get $0
    )
  
  (func $60 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    global.get $24
    i32.const 32
    i32.sub
    local.tee $3
    global.set $24
    block $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.eq
      local.get $0
      local.get $0
      f64.ne
      i32.or
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      block $block_0
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        i32.const 1
        i32.xor
        br_if $block_0
        local.get $3
        local.get $0
        f64.neg
        local.get $3
        local.get $3
        call $61
        local.get $3
        f64.load
        local.tee $0
        f64.const 0x1.0000000000000p-0
        f64.add
        local.get $0
        local.get $3
        f64.load offset=8
        f64.const 0x0.0000000000000p-1023
        f64.ne
        select
        f64.neg
        local.set $0
        br $block
      end ;; $block_0
      local.get $3
      i32.const 16
      i32.add
      local.get $0
      local.get $3
      local.get $3
      call $61
      local.get $3
      f64.load offset=16
      local.set $0
    end ;; $block
    local.get $3
    i32.const 32
    i32.add
    global.set $24
    local.get $0
    )
  
  (func $61 (type $9)
    (param $0 i32)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (local $4 i32)
    (local $5 f64)
    (local $6 i64)
    (local $7 i32)
    (local $8 f64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $4
    global.set $24
    block $block
      block $block_0
        block $block_1
          local.get $1
          f64.const 0x1.0000000000000p-0
          f64.lt
          i32.const 1
          i32.xor
          br_if $block_1
          local.get $1
          f64.const 0x0.0000000000000p-1023
          f64.lt
          i32.const 1
          i32.xor
          i32.eqz
          br_if $block
          block $block_2
            local.get $1
            f64.const 0x0.0000000000000p-1023
            f64.ne
            br_if $block_2
            local.get $1
            local.set $5
            br $block_0
          end ;; $block_2
          local.get $1
          local.set $5
          f64.const 0x0.0000000000000p-1023
          local.set $1
          br $block_0
        end ;; $block_1
        local.get $1
        i64.const -1
        i32.const 1075
        local.get $1
        i64.reinterpret_f64
        local.tee $6
        i64.const 52
        i64.shr_u
        i32.wrap_i64
        i32.const 2047
        i32.and
        local.tee $7
        i32.sub
        i64.extend_i32_u
        i64.shl
        i64.const -1
        local.get $7
        i32.const -1023
        i32.add
        i32.const 52
        i32.lt_u
        select
        local.get $6
        i64.and
        f64.reinterpret_i64
        local.tee $8
        f64.sub
        local.set $5
        local.get $8
        local.set $1
      end ;; $block_0
      local.get $0
      local.get $1
      f64.store
      local.get $0
      local.get $5
      f64.store offset=8
      local.get $4
      i32.const 16
      i32.add
      global.set $24
      return
    end ;; $block
    local.get $4
    local.get $1
    f64.neg
    local.get $0
    local.get $0
    call $61
    local.get $4
    f64.load
    local.set $1
    local.get $0
    local.get $4
    f64.load offset=8
    f64.neg
    f64.store offset=8
    local.get $0
    local.get $1
    f64.neg
    f64.store
    local.get $4
    i32.const 16
    i32.add
    global.set $24
    )
  
  (func $62 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    local.get $0
    f64.neg
    local.get $3
    local.get $3
    call $60
    f64.neg
    )
  
  (func $63 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $3
    global.set $24
    block $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.eq
      local.get $0
      local.get $0
      f64.ne
      i32.or
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $3
      local.get $0
      local.get $3
      local.get $3
      call $61
      local.get $3
      f64.load
      local.set $0
    end ;; $block
    local.get $3
    i32.const 16
    i32.add
    global.set $24
    local.get $0
    )
  
  (func $64 (type $12)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 i32)
    (local $5 f64)
    (local $6 f64)
    (local $7 f64)
    (local $8 i32)
    (local $9 i64)
    (local $10 i32)
    (local $11 i64)
    (local $12 i32)
    (local $13 i32)
    global.get $24
    i32.const 32
    i32.sub
    local.tee $4
    global.set $24
    f64.const 0x1.0000000000000p-0
    local.set $5
    block $block
      block $block_0
        block $block_1
          loop $loop
            local.get $0
            f64.const 0x1.0000000000000p-0
            f64.eq
            br_if $block
            local.get $1
            f64.const 0x0.0000000000000p-1023
            f64.eq
            br_if $block
            block $block_2
              local.get $1
              f64.const 0x1.0000000000000p-0
              f64.ne
              br_if $block_2
              local.get $0
              local.set $5
              br $block
            end ;; $block_2
            block $block_3
              local.get $0
              local.get $0
              f64.ne
              local.get $1
              local.get $1
              f64.ne
              i32.or
              i32.eqz
              br_if $block_3
              f64.const +nan:0x8000000000001
              local.set $5
              br $block
            end ;; $block_3
            block $block_4
              local.get $0
              f64.const 0x0.0000000000000p-1023
              f64.ne
              br_if $block_4
              block $block_5
                local.get $1
                f64.const 0x0.0000000000000p-1023
                f64.lt
                i32.const 1
                i32.xor
                br_if $block_5
                f64.const +inf
                local.set $5
                local.get $1
                call $65
                i32.const 1
                i32.and
                i32.eqz
                br_if $block
                local.get $0
                i64.reinterpret_f64
                i64.const -9223372036854775808
                i64.and
                i64.const 9218868437227405312
                i64.or
                f64.reinterpret_i64
                local.set $5
                br $block
              end ;; $block_5
              local.get $1
              f64.const 0x0.0000000000000p-1023
              f64.gt
              i32.const 1
              i32.xor
              br_if $block_0
              local.get $0
              f64.const 0x0.0000000000000p-1023
              local.get $1
              call $65
              i32.const 1
              i32.and
              select
              local.set $5
              br $block
            end ;; $block_4
            block $block_6
              block $block_7
                block $block_8
                  local.get $1
                  f64.const 0x1.fffffffffffffp+1023
                  f64.gt
                  br_if $block_8
                  local.get $1
                  f64.const -0x1.fffffffffffffp+1023
                  f64.lt
                  i32.eqz
                  br_if $block_6
                  local.get $0
                  f64.const -0x1.0000000000000p-0
                  f64.eq
                  br_if $block
                  f64.const 0x0.0000000000000p-1023
                  local.set $5
                  local.get $0
                  i64.reinterpret_f64
                  i64.const 9223372036854775807
                  i64.and
                  f64.reinterpret_i64
                  f64.const 0x1.0000000000000p-0
                  f64.lt
                  i32.const 1
                  i32.xor
                  i32.eqz
                  br_if $block_7
                  br $block
                end ;; $block_8
                local.get $0
                f64.const -0x1.0000000000000p-0
                f64.eq
                br_if $block
                f64.const 0x0.0000000000000p-1023
                local.set $5
                local.get $0
                i64.reinterpret_f64
                i64.const 9223372036854775807
                i64.and
                f64.reinterpret_i64
                f64.const 0x1.0000000000000p-0
                f64.lt
                br_if $block
              end ;; $block_7
              f64.const +inf
              local.set $5
              br $block
            end ;; $block_6
            block $block_9
              local.get $0
              f64.const 0x1.fffffffffffffp+1023
              f64.gt
              br_if $block_9
              local.get $0
              f64.const -0x1.fffffffffffffp+1023
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_1
            end ;; $block_9
            block $block_10
              local.get $0
              f64.const -0x1.fffffffffffffp+1023
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_10
              f64.const 0x1.0000000000000p-0
              local.get $0
              f64.div
              local.set $0
              local.get $1
              f64.neg
              local.set $1
              br $loop
            end ;; $block_10
          end ;; $loop
          f64.const 0x0.0000000000000p-1023
          local.set $5
          local.get $1
          f64.const 0x0.0000000000000p-1023
          f64.lt
          br_if $block
          f64.const +inf
          local.set $5
          local.get $1
          f64.const 0x0.0000000000000p-1023
          f64.gt
          i32.eqz
          br_if $block_0
          br $block
        end ;; $block_1
        block $block_11
          local.get $1
          f64.const 0x1.0000000000000p-1
          f64.ne
          br_if $block_11
          local.get $0
          local.get $4
          local.get $4
          call $38
          local.set $5
          br $block
        end ;; $block_11
        local.get $1
        f64.const -0x1.0000000000000p-1
        f64.ne
        br_if $block_0
        f64.const 0x1.0000000000000p-0
        local.get $0
        local.get $4
        local.get $4
        call $38
        f64.div
        local.set $5
        br $block
      end ;; $block_0
      local.get $4
      i32.const 16
      i32.add
      local.get $1
      i64.reinterpret_f64
      i64.const 9223372036854775807
      i64.and
      f64.reinterpret_i64
      local.get $4
      local.get $4
      call $61
      local.get $4
      f64.load offset=24
      local.set $6
      local.get $4
      f64.load offset=16
      local.set $7
      block $block_12
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        i32.const 1
        i32.xor
        br_if $block_12
        f64.const +nan:0x8000000000001
        local.set $5
        local.get $6
        f64.const 0x0.0000000000000p-1023
        f64.ne
        br_if $block
      end ;; $block_12
      block $block_13
        local.get $7
        f64.const 0x1.0000000000000p+63
        f64.ge
        i32.const 1
        i32.xor
        br_if $block_13
        f64.const 0x1.0000000000000p-0
        local.set $5
        local.get $0
        f64.const -0x1.0000000000000p-0
        f64.eq
        br_if $block
        local.get $4
        i32.const 32
        i32.add
        global.set $24
        f64.const +inf
        f64.const 0x0.0000000000000p-1023
        local.get $0
        i64.reinterpret_f64
        i64.const 9223372036854775807
        i64.and
        f64.reinterpret_i64
        f64.const 0x1.0000000000000p-0
        f64.lt
        local.get $1
        f64.const 0x0.0000000000000p-1023
        f64.gt
        i32.xor
        select
        return
      end ;; $block_13
      block $block_14
        block $block_15
          local.get $6
          f64.const 0x0.0000000000000p-1023
          f64.ne
          br_if $block_15
          f64.const 0x1.0000000000000p-0
          local.set $5
          br $block_14
        end ;; $block_15
        local.get $7
        f64.const 0x1.0000000000000p-0
        f64.add
        local.get $7
        local.get $6
        f64.const 0x1.0000000000000p-1
        f64.gt
        local.tee $8
        select
        local.set $7
        local.get $6
        f64.const -0x1.0000000000000p-0
        f64.add
        local.get $6
        local.get $8
        select
        local.get $0
        local.get $4
        local.get $4
        call $37
        f64.mul
        local.get $4
        local.get $4
        call $54
        local.set $5
      end ;; $block_14
      local.get $4
      local.get $0
      local.get $4
      local.get $4
      call $40
      i64.const 0
      i64.const 9223372036854775807
      i64.const -9223372036854775808
      local.get $7
      f64.const -0x1.0000000000000p+63
      f64.ge
      local.tee $8
      select
      local.get $7
      local.get $7
      f64.ne
      select
      local.set $9
      local.get $7
      f64.const 0x1.ffffffffffffep+62
      f64.le
      local.set $10
      block $block_16
        block $block_17
          local.get $7
          f64.abs
          f64.const 0x1.0000000000000p+63
          f64.lt
          i32.eqz
          br_if $block_17
          local.get $7
          i64.trunc_f64_s
          local.set $11
          br $block_16
        end ;; $block_17
        i64.const -9223372036854775808
        local.set $11
      end ;; $block_16
      local.get $11
      local.get $9
      local.get $10
      select
      local.get $9
      local.get $8
      select
      local.set $9
      local.get $4
      i32.load offset=8
      local.set $8
      local.get $4
      f64.load
      local.set $0
      i32.const 0
      local.set $10
      loop $loop_0
        block $block_18
          block $block_19
            local.get $9
            i64.eqz
            br_if $block_19
            local.get $8
            i32.const 4096
            i32.add
            i32.const 8193
            i32.lt_u
            br_if $block_18
            local.get $8
            local.get $10
            i32.add
            local.set $10
          end ;; $block_19
          block $block_20
            local.get $1
            f64.const 0x0.0000000000000p-1023
            f64.lt
            i32.const 1
            i32.xor
            br_if $block_20
            i32.const 0
            local.get $10
            i32.sub
            local.set $10
            f64.const 0x1.0000000000000p-0
            local.get $5
            f64.div
            local.set $5
          end ;; $block_20
          local.get $5
          local.get $10
          local.get $4
          local.get $4
          call $56
          local.set $5
          br $block
        end ;; $block_18
        local.get $5
        local.get $5
        local.get $0
        f64.mul
        local.get $9
        i64.const 1
        i64.and
        i64.eqz
        local.tee $12
        select
        local.set $5
        local.get $0
        local.get $0
        f64.mul
        local.tee $0
        local.get $0
        f64.add
        local.get $0
        local.get $0
        f64.const 0x1.0000000000000p-1
        f64.lt
        local.tee $13
        select
        local.set $0
        local.get $9
        i64.const 1
        i64.shr_s
        local.set $9
        i32.const 0
        local.get $8
        local.get $12
        select
        local.get $10
        i32.add
        local.set $10
        local.get $8
        i32.const 1
        i32.shl
        local.get $13
        i32.sub
        local.set $8
        br $loop_0
      end ;; $loop_0
    end ;; $block
    local.get $4
    i32.const 32
    i32.add
    global.set $24
    local.get $5
    )
  
  (func $65 (type $15)
    (param $0 f64)
    (result i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i64)
    (local $5 i64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $1
    global.set $24
    local.get $1
    local.get $0
    local.get $1
    local.get $1
    call $61
    block $block
      block $block_0
        local.get $1
        f64.load offset=8
        f64.const 0x0.0000000000000p-1023
        f64.eq
        br_if $block_0
        i32.const 0
        local.set $2
        br $block
      end ;; $block_0
      local.get $1
      f64.load
      local.tee $0
      f64.const 0x1.ffffffffffffep+62
      f64.le
      local.set $2
      local.get $0
      local.get $0
      f64.eq
      local.get $0
      f64.const -0x1.0000000000000p+63
      f64.ge
      local.tee $3
      i32.and
      i64.extend_i32_u
      local.set $4
      block $block_1
        block $block_2
          local.get $0
          f64.abs
          f64.const 0x1.0000000000000p+63
          f64.lt
          i32.eqz
          br_if $block_2
          local.get $0
          i64.trunc_f64_s
          local.set $5
          br $block_1
        end ;; $block_2
        i64.const -9223372036854775808
        local.set $5
      end ;; $block_1
      local.get $5
      local.get $4
      local.get $2
      select
      local.get $4
      local.get $3
      select
      i32.wrap_i64
      i32.const 1
      i32.and
      local.set $2
    end ;; $block
    local.get $1
    i32.const 16
    i32.add
    global.set $24
    local.get $2
    )
  
  (func $66 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    (local $4 f64)
    (local $5 i32)
    (local $6 i64)
    (local $7 f64)
    (local $8 i32)
    (local $9 i32)
    (local $10 i64)
    (local $11 f64)
    (local $12 f64)
    (local $13 f64)
    (local $14 f64)
    (local $15 f64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $3
    global.set $24
    block $block
      block $block_0
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.eq
        local.get $0
        local.get $0
        f64.ne
        i32.or
        i32.eqz
        br_if $block_0
        local.get $0
        local.set $4
        br $block
      end ;; $block_0
      f64.const +nan:0x8000000000001
      local.set $4
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      block $block_1
        block $block_2
          local.get $0
          f64.neg
          local.get $0
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          local.tee $5
          select
          local.tee $4
          f64.const 0x1.0000000000000p+29
          f64.ge
          i32.const 1
          i32.xor
          br_if $block_2
          local.get $3
          local.get $4
          call $67
          local.get $3
          f64.load offset=8
          local.set $4
          local.get $3
          i64.load
          local.set $6
          br $block_1
        end ;; $block_2
        i64.const -1
        i64.const 0
        local.get $4
        f64.const 0x1.45f306dc9c883p-0
        f64.mul
        local.tee $7
        f64.const 0x0.0000000000000p-1023
        f64.ge
        local.tee $8
        select
        local.set $6
        local.get $7
        f64.const 0x1.ffffffffffffep+63
        f64.le
        local.set $9
        block $block_3
          block $block_4
            local.get $7
            f64.const 0x1.0000000000000p+64
            f64.lt
            local.get $8
            i32.and
            i32.eqz
            br_if $block_4
            local.get $7
            i64.trunc_f64_u
            local.set $10
            br $block_3
          end ;; $block_4
          i64.const 0
          local.set $10
        end ;; $block_3
        local.get $4
        local.get $10
        local.get $6
        local.get $9
        select
        local.get $6
        local.get $8
        select
        local.tee $6
        f64.convert_i64_u
        local.tee $7
        local.get $7
        f64.const 0x1.0000000000000p-0
        f64.add
        local.get $6
        i64.const 1
        i64.and
        local.tee $10
        i64.eqz
        select
        local.tee $7
        f64.const -0x1.921fb40000000p-1
        f64.mul
        f64.add
        local.get $7
        f64.const -0x1.4442d00000000p-25
        f64.mul
        f64.add
        local.get $7
        f64.const -0x1.8469898cc5170p-49
        f64.mul
        f64.add
        local.set $4
        local.get $6
        local.get $10
        i64.add
        i64.const 7
        i64.and
        local.set $6
      end ;; $block_1
      block $block_5
        local.get $6
        i64.const 4
        i64.lt_u
        br_if $block_5
        local.get $6
        i64.const -4
        i64.add
        local.set $6
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        i32.const 1
        i32.xor
        local.set $5
      end ;; $block_5
      local.get $4
      local.get $4
      f64.mul
      local.set $0
      block $block_6
        block $block_7
          local.get $6
          i64.const -1
          i64.add
          i64.const 1
          i64.gt_u
          br_if $block_7
          local.get $0
          f64.const -0x1.8fa49a0861a9bp-37
          f64.mul
          f64.const 0x1.1ee9d7b4e3f05p-29
          f64.add
          local.set $7
          local.get $0
          f64.const -0x1.0000000000000p-1
          f64.mul
          f64.const 0x1.0000000000000p-0
          f64.add
          local.set $11
          f64.const 0x1.555555555554bp-5
          local.set $12
          f64.const -0x1.6c16c16c14f91p-10
          local.set $13
          f64.const 0x1.a01a019c844f5p-16
          local.set $14
          f64.const -0x1.27e4f7eac4bc6p-22
          local.set $15
          local.get $0
          local.set $4
          br $block_6
        end ;; $block_7
        local.get $0
        f64.const 0x1.5d8fd1fd19ccdp-33
        f64.mul
        f64.const -0x1.ae5e5a9291f5dp-26
        f64.add
        local.set $7
        f64.const -0x1.5555555555548p-3
        local.set $12
        f64.const 0x1.111111110f7d0p-7
        local.set $13
        f64.const -0x1.a01a019bfdf03p-13
        local.set $14
        f64.const 0x1.71de3567d48a1p-19
        local.set $15
        local.get $4
        local.set $11
      end ;; $block_6
      local.get $11
      local.get $4
      local.get $0
      f64.mul
      local.get $0
      local.get $0
      local.get $0
      local.get $0
      local.get $7
      f64.mul
      local.get $15
      f64.add
      f64.mul
      local.get $14
      f64.add
      f64.mul
      local.get $13
      f64.add
      f64.mul
      local.get $12
      f64.add
      f64.mul
      f64.add
      local.tee $0
      f64.neg
      local.get $0
      local.get $5
      select
      local.set $4
    end ;; $block
    local.get $3
    i32.const 16
    i32.add
    global.set $24
    local.get $4
    )
  
  (func $67 (type $10)
    (param $0 i32)
    (param $1 f64)
    (local $2 i32)
    (local $3 i64)
    (local $4 i32)
    (local $5 i32)
    (local $6 i64)
    (local $7 i64)
    (local $8 i64)
    (local $9 i32)
    (local $10 i32)
    (local $11 i32)
    global.get $24
    i32.const 32
    i32.sub
    local.tee $2
    global.set $24
    block $block
      block $block_0
        local.get $1
        f64.const 0x1.921fb54442d18p-1
        f64.lt
        i32.const 1
        i32.xor
        i32.eqz
        br_if $block_0
        local.get $1
        i64.reinterpret_f64
        local.tee $3
        i64.const 52
        i64.shr_u
        i32.wrap_i64
        i32.const 2047
        i32.and
        i32.const -1014
        i32.add
        local.tee $4
        i32.const 1280
        i32.ge_u
        br_if $block
        local.get $2
        i32.const 16
        i32.add
        local.get $4
        i32.const 6
        i32.shr_u
        i32.const 3
        i32.shl
        local.tee $5
        i32.const 65848
        i32.add
        i64.load
        i32.const 64
        local.get $4
        i32.const 63
        i32.and
        local.tee $4
        i32.sub
        i64.extend_i32_u
        local.tee $6
        i64.shr_u
        i64.const 0
        local.get $4
        select
        local.get $5
        i32.const 65840
        i32.add
        i64.load
        local.tee $7
        local.get $4
        i64.extend_i32_u
        local.tee $8
        i64.shl
        i64.or
        local.get $3
        i64.const -9218868437227405313
        i64.and
        i64.const 4503599627370496
        i64.or
        local.tee $3
        call $35
        local.get $2
        local.get $7
        local.get $6
        i64.shr_u
        i64.const 0
        local.get $4
        select
        local.get $5
        i32.const 65832
        i32.add
        i64.load
        local.tee $7
        local.get $8
        i64.shl
        i64.or
        local.get $3
        call $35
        local.get $2
        i64.load
        local.get $7
        local.get $6
        i64.shr_u
        i64.const 0
        local.get $4
        select
        local.get $5
        i32.const 65824
        i32.add
        i64.load
        local.get $8
        i64.shl
        i64.or
        local.get $3
        i64.mul
        i64.add
        local.get $2
        i64.load offset=8
        local.tee $6
        local.get $2
        i64.load offset=16
        local.tee $8
        i64.add
        local.tee $3
        i64.const -1
        i64.xor
        local.get $6
        local.get $8
        i64.or
        i64.and
        local.get $6
        local.get $8
        i64.and
        i64.or
        i64.const 63
        i64.shr_u
        i64.add
        local.tee $6
        i64.const 29
        i64.shr_u
        i64.const 4294967295
        i64.and
        local.get $6
        i64.const 3
        i64.shl
        local.get $3
        i64.const 61
        i64.shr_u
        i64.or
        local.tee $8
        local.get $8
        i64.const 4294967295
        i64.gt_u
        local.tee $4
        select
        local.tee $7
        i64.const 16
        i64.shr_u
        local.get $7
        local.get $7
        i64.const 65535
        i64.gt_u
        local.tee $5
        select
        local.tee $7
        i64.const 8
        i64.shr_u
        local.get $7
        local.get $7
        i64.const 255
        i64.gt_u
        local.tee $9
        select
        i32.wrap_i64
        local.tee $10
        i32.const 256
        i32.ge_u
        br_if $block
        local.get $0
        local.get $6
        i64.const 61
        i64.shr_u
        local.tee $7
        local.get $7
        i64.const 1
        i64.add
        i64.const 7
        i64.and
        local.get $6
        i64.const 2305843009213693952
        i64.and
        i64.eqz
        local.tee $11
        select
        i64.store
        local.get $0
        i64.const 0
        local.get $8
        i32.const 65
        local.get $4
        i32.const 5
        i32.shl
        local.tee $4
        i32.const 16
        i32.or
        local.get $4
        local.get $5
        select
        local.tee $4
        i32.const 8
        i32.or
        local.get $4
        local.get $9
        select
        local.get $10
        i32.const 65568
        i32.add
        i32.load8_u
        i32.add
        local.tee $4
        i32.sub
        local.tee $5
        i64.extend_i32_u
        i64.shl
        local.get $5
        i32.const 63
        i32.gt_u
        select
        i64.const 0
        local.get $3
        i32.const 63
        i32.const 64
        local.get $4
        i32.sub
        local.tee $4
        i32.sub
        i64.extend_i32_u
        i64.shr_u
        local.get $4
        i32.const 63
        i32.gt_u
        select
        i64.or
        i64.const 12
        i64.shr_u
        i32.const 1022
        local.get $4
        i32.sub
        i64.extend_i32_u
        i64.const 52
        i64.shl
        i64.or
        f64.reinterpret_i64
        local.tee $1
        local.get $1
        f64.const -0x1.0000000000000p-0
        f64.add
        local.get $11
        select
        f64.const 0x1.921fb54442d18p-1
        f64.mul
        f64.store offset=8
        local.get $2
        i32.const 32
        i32.add
        global.set $24
        return
      end ;; $block_0
      local.get $0
      local.get $1
      f64.store offset=8
      local.get $0
      i64.const 0
      i64.store
      local.get $2
      i32.const 32
      i32.add
      global.set $24
      return
    end ;; $block
    call $68
    unreachable
    )
  
  (func $68 (type $5)
    i32.const 66117
    i32.const 18
    call $100
    )
  
  (func $69 (type $12)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 f64)
    (local $5 i32)
    (local $6 f64)
    f64.const +inf
    local.set $4
    block $block
      local.get $1
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $1
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      block $block_0
        local.get $0
        local.get $0
        f64.ne
        local.get $1
        local.get $1
        f64.ne
        i32.or
        i32.eqz
        br_if $block_0
        f64.const +nan:0x8000000000001
        return
      end ;; $block_0
      f64.const 0x0.0000000000000p-1023
      local.set $4
      local.get $1
      i64.reinterpret_f64
      i64.const 9223372036854775807
      i64.and
      f64.reinterpret_i64
      local.tee $1
      local.get $0
      i64.reinterpret_f64
      i64.const 9223372036854775807
      i64.and
      f64.reinterpret_i64
      local.tee $0
      local.get $0
      local.get $1
      f64.lt
      local.tee $5
      select
      local.tee $6
      f64.const 0x0.0000000000000p-1023
      f64.eq
      br_if $block
      local.get $6
      local.get $0
      local.get $1
      local.get $5
      select
      local.get $6
      f64.div
      local.tee $1
      local.get $1
      f64.mul
      f64.const 0x1.0000000000000p-0
      f64.add
      local.get $5
      local.get $5
      call $38
      f64.mul
      local.set $4
    end ;; $block
    local.get $4
    )
  
  (func $70 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    (local $4 f64)
    (local $5 i64)
    (local $6 i32)
    (local $7 i32)
    (local $8 i64)
    (local $9 f64)
    (local $10 f64)
    (local $11 f64)
    (local $12 f64)
    (local $13 f64)
    (local $14 f64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $3
    global.set $24
    f64.const +nan:0x8000000000001
    local.set $4
    block $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      local.get $0
      f64.ne
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      block $block_0
        block $block_1
          local.get $0
          i64.reinterpret_f64
          i64.const 9223372036854775807
          i64.and
          f64.reinterpret_i64
          local.tee $0
          f64.const 0x1.0000000000000p+29
          f64.ge
          i32.const 1
          i32.xor
          br_if $block_1
          local.get $3
          local.get $0
          call $67
          local.get $3
          f64.load offset=8
          local.set $4
          local.get $3
          i64.load
          local.set $5
          br $block_0
        end ;; $block_1
        i64.const -1
        i64.const 0
        local.get $0
        f64.const 0x1.45f306dc9c883p-0
        f64.mul
        local.tee $4
        f64.const 0x0.0000000000000p-1023
        f64.ge
        local.tee $6
        select
        local.set $5
        local.get $4
        f64.const 0x1.ffffffffffffep+63
        f64.le
        local.set $7
        block $block_2
          block $block_3
            local.get $4
            f64.const 0x1.0000000000000p+64
            f64.lt
            local.get $6
            i32.and
            i32.eqz
            br_if $block_3
            local.get $4
            i64.trunc_f64_u
            local.set $8
            br $block_2
          end ;; $block_3
          i64.const 0
          local.set $8
        end ;; $block_2
        local.get $0
        local.get $8
        local.get $5
        local.get $7
        select
        local.get $5
        local.get $6
        select
        local.tee $5
        f64.convert_i64_u
        local.tee $4
        local.get $4
        f64.const 0x1.0000000000000p-0
        f64.add
        local.get $5
        i64.const 1
        i64.and
        local.tee $8
        i64.eqz
        select
        local.tee $4
        f64.const -0x1.921fb40000000p-1
        f64.mul
        f64.add
        local.get $4
        f64.const -0x1.4442d00000000p-25
        f64.mul
        f64.add
        local.get $4
        f64.const -0x1.8469898cc5170p-49
        f64.mul
        f64.add
        local.set $4
        local.get $5
        local.get $8
        i64.add
        i64.const 7
        i64.and
        local.set $5
      end ;; $block_0
      local.get $4
      local.get $4
      f64.mul
      local.set $0
      block $block_4
        block $block_5
          local.get $5
          i64.const -4
          i64.add
          local.get $5
          local.get $5
          i64.const 3
          i64.gt_u
          select
          local.tee $8
          i64.const -1
          i64.add
          i64.const 1
          i64.gt_u
          br_if $block_5
          local.get $0
          f64.const 0x1.5d8fd1fd19ccdp-33
          f64.mul
          f64.const -0x1.ae5e5a9291f5dp-26
          f64.add
          local.set $9
          f64.const -0x1.5555555555548p-3
          local.set $10
          f64.const 0x1.111111110f7d0p-7
          local.set $11
          f64.const -0x1.a01a019bfdf03p-13
          local.set $12
          f64.const 0x1.71de3567d48a1p-19
          local.set $13
          local.get $4
          local.set $14
          br $block_4
        end ;; $block_5
        local.get $0
        f64.const -0x1.8fa49a0861a9bp-37
        f64.mul
        f64.const 0x1.1ee9d7b4e3f05p-29
        f64.add
        local.set $9
        local.get $0
        f64.const -0x1.0000000000000p-1
        f64.mul
        f64.const 0x1.0000000000000p-0
        f64.add
        local.set $14
        f64.const 0x1.555555555554bp-5
        local.set $10
        f64.const -0x1.6c16c16c14f91p-10
        local.set $11
        f64.const 0x1.a01a019c844f5p-16
        local.set $12
        f64.const -0x1.27e4f7eac4bc6p-22
        local.set $13
        local.get $0
        local.set $4
      end ;; $block_4
      local.get $14
      local.get $4
      local.get $0
      f64.mul
      local.get $0
      local.get $0
      local.get $0
      local.get $0
      local.get $9
      f64.mul
      local.get $13
      f64.add
      f64.mul
      local.get $12
      f64.add
      f64.mul
      local.get $11
      f64.add
      f64.mul
      local.get $10
      f64.add
      f64.mul
      f64.add
      local.tee $0
      f64.neg
      local.get $0
      local.get $5
      i64.const 3
      i64.gt_u
      local.get $8
      i64.const 1
      i64.gt_u
      i32.xor
      select
      local.set $4
    end ;; $block
    local.get $3
    i32.const 16
    i32.add
    global.set $24
    local.get $4
    )
  
  (func $71 (type $12)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 i32)
    (local $5 f64)
    (local $6 i32)
    (local $7 f64)
    global.get $24
    i32.const 32
    i32.sub
    local.tee $4
    global.set $24
    f64.const +nan:0x8000000000001
    local.set $5
    block $block
      local.get $1
      local.get $1
      f64.ne
      local.get $0
      local.get $0
      f64.ne
      i32.or
      br_if $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $1
      f64.const 0x0.0000000000000p-1023
      f64.eq
      br_if $block
      local.get $4
      i32.const 16
      i32.add
      local.get $1
      i64.reinterpret_f64
      i64.const 9223372036854775807
      i64.and
      f64.reinterpret_i64
      local.tee $5
      local.get $4
      local.get $4
      call $40
      local.get $0
      f64.neg
      local.get $0
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      select
      local.set $1
      local.get $4
      i32.load offset=24
      local.set $6
      local.get $4
      f64.load offset=16
      local.set $7
      block $block_0
        loop $loop
          local.get $1
          local.get $5
          f64.ge
          i32.const 1
          i32.xor
          br_if $block_0
          local.get $4
          local.get $1
          local.get $4
          local.get $4
          call $40
          local.get $1
          local.get $5
          local.get $4
          i32.load offset=8
          local.get $6
          i32.sub
          local.get $4
          f64.load
          local.get $7
          f64.lt
          i32.sub
          local.get $4
          local.get $4
          call $56
          f64.sub
          local.set $1
          br $loop
        end ;; $loop
      end ;; $block_0
      local.get $1
      f64.neg
      local.get $1
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      select
      local.set $5
    end ;; $block
    local.get $4
    i32.const 32
    i32.add
    global.set $24
    local.get $5
    )
  
  (func $72 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    local.get $0
    local.get $3
    local.get $3
    call $37
    f64.const 0x1.bcb7b1526e50ep-2
    f64.mul
    )
  
  (func $73 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    (local $4 i32)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $3
    global.set $24
    local.get $3
    local.get $0
    local.get $3
    local.get $3
    call $40
    local.get $3
    i32.load offset=8
    local.set $4
    block $block
      block $block_0
        local.get $3
        f64.load
        local.tee $0
        f64.const 0x1.0000000000000p-1
        f64.ne
        br_if $block_0
        local.get $4
        i32.const -1
        i32.add
        f64.convert_i32_s
        local.set $0
        br $block
      end ;; $block_0
      local.get $0
      local.get $3
      local.get $3
      call $37
      f64.const 0x1.71547652b82fep-0
      f64.mul
      local.get $4
      f64.convert_i32_s
      f64.add
      local.set $0
    end ;; $block
    local.get $3
    i32.const 16
    i32.add
    global.set $24
    local.get $0
    )
  
  (func $74 (type $12)
    (param $0 f64)
    (param $1 f64)
    (param $2 i32)
    (param $3 i32)
    (result f64)
    (local $4 f64)
    (local $5 i32)
    (local $6 f64)
    f64.const +nan:0x8000000000001
    local.set $4
    block $block
      block $block_0
        local.get $0
        f64.const 0x1.fffffffffffffp+1023
        f64.gt
        br_if $block_0
        local.get $0
        local.get $0
        f64.ne
        local.get $1
        local.get $1
        f64.ne
        i32.or
        br_if $block_0
        local.get $0
        f64.const -0x1.fffffffffffffp+1023
        f64.lt
        br_if $block_0
        local.get $1
        f64.const 0x0.0000000000000p-1023
        f64.eq
        br_if $block_0
        local.get $1
        f64.const 0x1.fffffffffffffp+1023
        f64.gt
        br_if $block
        local.get $1
        f64.const -0x1.fffffffffffffp+1023
        f64.lt
        br_if $block
        block $block_1
          local.get $0
          f64.neg
          local.get $0
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          local.tee $5
          select
          local.tee $4
          local.get $1
          f64.neg
          local.get $1
          local.get $1
          f64.const 0x0.0000000000000p-1023
          f64.lt
          select
          local.tee $1
          f64.ne
          br_if $block_1
          f64.const -0x0.0000000000000p-1023
          f64.const 0x0.0000000000000p-1023
          local.get $5
          select
          return
        end ;; $block_1
        block $block_2
          local.get $1
          f64.const 0x1.fffffffffffffp+1022
          f64.le
          i32.const 1
          i32.xor
          br_if $block_2
          local.get $4
          local.get $1
          local.get $1
          f64.add
          local.get $5
          local.get $5
          call $71
          local.set $4
        end ;; $block_2
        block $block_3
          block $block_4
            block $block_5
              local.get $1
              f64.const 0x1.0000000000000p-1021
              f64.lt
              i32.const 1
              i32.xor
              br_if $block_5
              local.get $4
              local.get $4
              f64.add
              local.get $1
              f64.gt
              i32.const 1
              i32.xor
              br_if $block_3
              local.get $4
              local.get $1
              f64.sub
              local.tee $4
              local.get $4
              f64.add
              local.get $1
              f64.ge
              i32.const 1
              i32.xor
              br_if $block_3
              br $block_4
            end ;; $block_5
            local.get $4
            local.get $1
            f64.const 0x1.0000000000000p-1
            f64.mul
            local.tee $6
            f64.gt
            i32.const 1
            i32.xor
            br_if $block_3
            local.get $4
            local.get $1
            f64.sub
            local.tee $4
            local.get $6
            f64.ge
            i32.const 1
            i32.xor
            br_if $block_3
          end ;; $block_4
          local.get $4
          local.get $1
          f64.sub
          local.set $4
        end ;; $block_3
        local.get $4
        f64.neg
        local.get $4
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        select
        local.set $4
      end ;; $block_0
      local.get $4
      return
    end ;; $block
    local.get $0
    )
  
  (func $75 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    (local $5 f64)
    block $block
      block $block_0
        local.get $0
        f64.neg
        local.get $0
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        select
        local.tee $3
        f64.const 0x1.5000000000000p+4
        f64.gt
        i32.const 1
        i32.xor
        br_if $block_0
        local.get $3
        local.get $4
        local.get $4
        call $54
        f64.const 0x1.0000000000000p-1
        f64.mul
        local.set $3
        br $block
      end ;; $block_0
      block $block_1
        local.get $3
        f64.const 0x1.0000000000000p-1
        f64.gt
        i32.const 1
        i32.xor
        br_if $block_1
        local.get $3
        local.get $4
        local.get $4
        call $54
        local.tee $3
        f64.const -0x1.0000000000000p-0
        local.get $3
        f64.div
        f64.add
        f64.const 0x1.0000000000000p-1
        f64.mul
        local.set $3
        br $block
      end ;; $block_1
      local.get $3
      local.get $3
      local.get $3
      f64.mul
      local.tee $5
      local.get $5
      local.get $5
      f64.const -0x1.a4e3de8540779p+4
      f64.mul
      f64.const -0x1.69c6c36da2dfbp+11
      f64.add
      f64.mul
      f64.const -0x1.5f38b8605d22dp+16
      f64.add
      f64.mul
      f64.const -0x1.33fdeba64bb4fp+19
      f64.add
      f64.mul
      local.get $5
      local.get $5
      local.get $5
      f64.const -0x1.5b5b9fcd003bbp+7
      f64.add
      f64.mul
      f64.const 0x1.db7963eae91e1p+13
      f64.add
      f64.mul
      f64.const -0x1.33fdeba64bb4fp+19
      f64.add
      f64.div
      local.set $3
    end ;; $block
    local.get $3
    f64.neg
    local.get $3
    local.get $0
    f64.const 0x0.0000000000000p-1023
    f64.lt
    select
    )
  
  (func $76 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    local.get $0
    i64.reinterpret_f64
    i64.const 9223372036854775807
    i64.and
    f64.reinterpret_i64
    local.tee $3
    local.get $4
    local.get $4
    call $54
    local.set $0
    block $block
      local.get $3
      f64.const 0x1.5000000000000p+4
      f64.gt
      br_if $block
      local.get $0
      f64.const 0x1.0000000000000p-0
      local.get $0
      f64.div
      f64.add
      local.set $0
    end ;; $block
    local.get $0
    f64.const 0x1.0000000000000p-1
    f64.mul
    )
  
  (func $77 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 i32)
    (local $4 f64)
    (local $5 i64)
    (local $6 f64)
    (local $7 i32)
    (local $8 i32)
    (local $9 i64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $3
    global.set $24
    block $block
      block $block_0
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.eq
        local.get $0
        local.get $0
        f64.ne
        i32.or
        i32.eqz
        br_if $block_0
        local.get $0
        local.set $4
        br $block
      end ;; $block_0
      f64.const +nan:0x8000000000001
      local.set $4
      local.get $0
      f64.const 0x1.fffffffffffffp+1023
      f64.gt
      br_if $block
      local.get $0
      f64.const -0x1.fffffffffffffp+1023
      f64.lt
      br_if $block
      block $block_1
        block $block_2
          local.get $0
          f64.neg
          local.get $0
          local.get $0
          f64.const 0x0.0000000000000p-1023
          f64.lt
          select
          local.tee $4
          f64.const 0x1.0000000000000p+29
          f64.ge
          i32.const 1
          i32.xor
          br_if $block_2
          local.get $3
          local.get $4
          call $67
          local.get $3
          f64.load offset=8
          local.set $4
          local.get $3
          i64.load
          local.set $5
          br $block_1
        end ;; $block_2
        i64.const -1
        i64.const 0
        local.get $4
        f64.const 0x1.45f306dc9c883p-0
        f64.mul
        local.tee $6
        f64.const 0x0.0000000000000p-1023
        f64.ge
        local.tee $7
        select
        local.set $5
        local.get $6
        f64.const 0x1.ffffffffffffep+63
        f64.le
        local.set $8
        block $block_3
          block $block_4
            local.get $6
            f64.const 0x1.0000000000000p+64
            f64.lt
            local.get $7
            i32.and
            i32.eqz
            br_if $block_4
            local.get $6
            i64.trunc_f64_u
            local.set $9
            br $block_3
          end ;; $block_4
          i64.const 0
          local.set $9
        end ;; $block_3
        local.get $4
        local.get $9
        local.get $5
        local.get $8
        select
        local.get $5
        local.get $7
        select
        local.tee $5
        f64.convert_i64_u
        local.tee $6
        local.get $6
        f64.const 0x1.0000000000000p-0
        f64.add
        local.get $5
        i64.const 1
        i64.and
        local.tee $9
        i64.eqz
        select
        local.tee $6
        f64.const -0x1.921fb40000000p-1
        f64.mul
        f64.add
        local.get $6
        f64.const -0x1.4442d00000000p-25
        f64.mul
        f64.add
        local.get $6
        f64.const -0x1.8469898cc5170p-49
        f64.mul
        f64.add
        local.set $4
        local.get $5
        local.get $9
        i64.add
        local.set $5
      end ;; $block_1
      block $block_5
        local.get $4
        local.get $4
        f64.mul
        local.tee $6
        f64.const 0x1.6849b86a12b9bp-47
        f64.gt
        i32.const 1
        i32.xor
        br_if $block_5
        local.get $4
        local.get $4
        local.get $6
        local.get $6
        local.get $6
        f64.const -0x1.992d8d24f3f38p+13
        f64.mul
        f64.const 0x1.199eca5fc9dddp+20
        f64.add
        f64.mul
        f64.const -0x1.11fead3299176p+24
        f64.add
        f64.mul
        local.get $6
        local.get $6
        local.get $6
        local.get $6
        f64.const 0x1.ab8a5eeb36572p+13
        f64.add
        f64.mul
        f64.const -0x1.427bc582abc96p+20
        f64.add
        f64.mul
        f64.const 0x1.7d98fc2ead8efp+24
        f64.add
        f64.mul
        f64.const -0x1.9afe03cbe5a31p+25
        f64.add
        f64.div
        f64.mul
        f64.add
        local.set $4
      end ;; $block_5
      local.get $4
      f64.const -0x1.0000000000000p-0
      local.get $4
      f64.div
      local.get $5
      i64.const 2
      i64.and
      i64.eqz
      select
      local.tee $4
      f64.neg
      local.get $4
      local.get $0
      f64.const 0x0.0000000000000p-1023
      f64.lt
      select
      local.set $4
    end ;; $block
    local.get $3
    i32.const 16
    i32.add
    global.set $24
    local.get $4
    )
  
  (func $78 (type $8)
    (param $0 f64)
    (param $1 i32)
    (param $2 i32)
    (result f64)
    (local $3 f64)
    (local $4 i32)
    block $block
      block $block_0
        local.get $0
        i64.reinterpret_f64
        i64.const 9223372036854775807
        i64.and
        f64.reinterpret_i64
        local.tee $3
        f64.const 0x1.601e678fc457bp+5
        f64.gt
        i32.const 1
        i32.xor
        br_if $block_0
        f64.const -0x1.0000000000000p-0
        local.set $3
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        br_if $block
        f64.const 0x1.0000000000000p-0
        return
      end ;; $block_0
      block $block_1
        local.get $3
        f64.const 0x1.4000000000000p-1
        f64.ge
        i32.const 1
        i32.xor
        br_if $block_1
        f64.const -0x1.0000000000000p+1
        local.get $3
        local.get $3
        f64.add
        local.get $4
        local.get $4
        call $54
        f64.const 0x1.0000000000000p-0
        f64.add
        f64.div
        f64.const 0x1.0000000000000p-0
        f64.add
        local.set $3
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.lt
        i32.const 1
        i32.xor
        br_if $block
        local.get $3
        f64.neg
        return
      end ;; $block_1
      block $block_2
        local.get $0
        f64.const 0x0.0000000000000p-1023
        f64.ne
        br_if $block_2
        local.get $0
        return
      end ;; $block_2
      local.get $0
      local.get $0
      f64.mul
      local.tee $3
      local.get $0
      f64.mul
      local.get $3
      local.get $3
      f64.const -0x1.edc5baafd6f4bp-1
      f64.mul
      f64.const -0x1.8d26a0e26682dp+6
      f64.add
      f64.mul
      f64.const -0x1.93ac030580563p+10
      f64.add
      f64.mul
      local.get $3
      local.get $3
      local.get $3
      f64.const 0x1.c33f28a581b86p+6
      f64.add
      f64.mul
      f64.const 0x1.176fa0e5535fap+11
      f64.add
      f64.mul
      f64.const 0x1.2ec102442040cp+12
      f64.add
      f64.div
      local.get $0
      f64.add
      local.set $3
    end ;; $block
    local.get $3
    )
  
  (func $79 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    block $block
      block $block_0
        local.get $0
        i32.eqz
        br_if $block_0
        local.get $0
        local.get $1
        call $80
        local.tee $0
        i32.eqz
        br_if $block
        local.get $0
        i32.const 4
        i32.add
        return
      end ;; $block_0
      i32.const 65984
      call $81
      unreachable
    end ;; $block
    i32.const 0
    )
  
  (func $80 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i64)
    (local $8 i64)
    block $block
      block $block_0
        block $block_1
          local.get $1
          i32.const 1073741821
          i32.ge_u
          br_if $block_1
          block $block_2
            local.get $0
            local.get $1
            i32.const 19
            i32.add
            i32.const -16
            i32.and
            i32.const -4
            i32.add
            i32.const 12
            local.get $1
            i32.const 12
            i32.gt_u
            select
            local.tee $2
            call $82
            local.tee $1
            br_if $block_2
            local.get $0
            i32.eqz
            br_if $block_0
            local.get $0
            i32.const 20
            i32.add
            i32.load
            i32.eqz
            br_if $block
            i32.const 0
            local.set $1
            local.get $2
            local.set $3
            block $block_3
              local.get $2
              i32.const 536870909
              i32.gt_u
              br_if $block_3
              i32.const -1
              i32.const -1
              i32.const 27
              local.get $2
              call $83
              i32.sub
              local.tee $1
              i32.shl
              i32.const -1
              i32.xor
              local.get $1
              i32.const 31
              i32.gt_u
              select
              local.get $2
              i32.add
              local.set $3
              local.get $0
              i32.load offset=20
              i32.eqz
              local.set $1
            end ;; $block_3
            local.get $1
            br_if $block_0
            local.get $0
            i32.load
            i32.load offset=1568
            local.set $4
            local.get $0
            i32.const 80
            i32.add
            i32.load
            local.set $5
            memory.size
            local.set $1
            local.get $3
            i32.const 0
            i32.const 4
            local.get $5
            i32.const 16
            i32.shl
            local.get $4
            i32.const 4
            i32.ne
            i32.sub
            local.tee $4
            i32.shl
            local.get $4
            i32.const 31
            i32.gt_u
            select
            i32.add
            i32.const 65535
            i32.add
            i32.const 16
            i32.shr_s
            local.tee $6
            memory.grow
            drop
            i32.const 0
            local.get $1
            i32.const 16
            i32.shl
            local.get $1
            memory.size
            local.tee $3
            i32.eq
            local.tee $4
            select
            local.tee $1
            i32.const 0
            local.get $1
            select
            local.tee $5
            i32.eqz
            br_if $block
            i32.const 0
            local.get $3
            i32.const 16
            i32.shl
            local.get $4
            select
            i32.const 0
            local.get $1
            select
            local.tee $3
            i32.eqz
            br_if $block
            block $block_4
              local.get $6
              i32.const 0
              local.get $1
              select
              local.tee $1
              br_if $block_4
              local.get $3
              local.get $5
              i32.sub
              local.tee $1
              i32.const 16
              i32.shr_u
              local.get $1
              i32.const 65535
              i32.and
              i32.const 0
              i32.ne
              i32.add
              local.set $1
            end ;; $block_4
            local.get $0
            local.get $3
            i32.store offset=8
            local.get $0
            local.get $0
            i32.load offset=80
            local.get $1
            i32.add
            i32.store offset=80
            local.get $0
            local.get $5
            local.get $3
            call $84
            local.get $0
            local.get $2
            call $82
            local.tee $1
            i32.eqz
            br_if $block
          end ;; $block_2
          local.get $0
          local.get $1
          call $85
          block $block_5
            block $block_6
              local.get $1
              i32.load
              local.tee $3
              i32.const -4
              i32.and
              local.tee $4
              local.get $2
              i32.sub
              local.tee $5
              i32.const 16
              i32.lt_u
              br_if $block_6
              local.get $1
              local.get $3
              i32.const 2
              i32.and
              local.get $2
              i32.or
              i32.store
              local.get $2
              local.get $1
              i32.add
              i32.const 4
              i32.add
              local.tee $2
              i32.eqz
              br_if $block_0
              local.get $2
              local.get $5
              i32.const -4
              i32.add
              i32.const 1
              i32.or
              i32.store
              local.get $0
              local.get $2
              call $86
              br $block_5
            end ;; $block_6
            local.get $1
            local.get $3
            i32.const -2
            i32.and
            i32.store
            local.get $1
            local.get $4
            i32.add
            i32.const 4
            i32.add
            local.tee $2
            i32.eqz
            br_if $block_0
            local.get $2
            local.get $2
            i32.load
            i32.const -3
            i32.and
            i32.store
          end ;; $block_5
          local.get $0
          i32.eqz
          br_if $block_0
          local.get $0
          i32.const 40
          i32.add
          local.tee $2
          local.get $2
          i64.load
          local.get $1
          i64.load32_u
          i64.const 4294967292
          i64.and
          local.tee $7
          i64.add
          local.tee $8
          i64.store
          block $block_7
            local.get $8
            local.get $0
            i32.const 48
            i32.add
            i64.load
            i64.le_s
            br_if $block_7
            local.get $0
            local.get $8
            i64.store offset=48
          end ;; $block_7
          local.get $0
          i32.const 56
          i32.add
          local.tee $2
          local.get $2
          i64.load
          local.get $7
          i64.sub
          i64.store
          local.get $0
          i32.const 64
          i32.add
          local.tee $0
          local.get $0
          i32.load
          i32.const 1
          i32.add
          i32.store
          local.get $1
          return
        end ;; $block_1
        i32.const 66216
        call $81
        unreachable
      end ;; $block_0
      call $87
      unreachable
    end ;; $block
    i32.const 0
    )
  
  (func $81 (type $0)
    (param $0 i32)
    i32.const 66185
    i32.const 7
    call $101
    local.get $0
    call $111
    call $102
    unreachable
    unreachable
    )
  
  (func $82 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    (local $2 i32)
    (local $3 i32)
    block $block
      local.get $0
      i32.eqz
      br_if $block
      local.get $0
      i32.load
      local.set $0
      block $block_0
        block $block_1
          local.get $1
          i32.const 256
          i32.lt_u
          br_if $block_1
          block $block_2
            local.get $1
            i32.const 536870909
            i32.gt_u
            br_if $block_2
            local.get $1
            i32.const 0
            i32.const 1
            i32.const 27
            local.get $1
            call $83
            i32.sub
            local.tee $2
            i32.shl
            local.get $2
            i32.const 31
            i32.gt_u
            select
            i32.add
            i32.const -1
            i32.add
            local.set $1
          end ;; $block_2
          i32.const 16
          local.get $1
          i32.const 27
          local.get $1
          call $83
          local.tee $3
          i32.sub
          local.tee $2
          i32.shr_u
          i32.const 16
          i32.xor
          local.get $2
          i32.const 31
          i32.gt_u
          select
          local.set $2
          i32.const 24
          local.get $3
          i32.sub
          local.set $1
          br $block_0
        end ;; $block_1
        local.get $1
        i32.const 4
        i32.shr_u
        local.set $2
        i32.const 0
        local.set $1
      end ;; $block_0
      block $block_3
        block $block_4
          i32.const 0
          i32.const -1
          local.get $2
          i32.shl
          local.get $2
          i32.const 31
          i32.gt_u
          select
          local.get $1
          i32.const 2
          i32.shl
          local.get $0
          i32.const 4
          i32.add
          local.tee $3
          i32.add
          i32.load
          i32.and
          local.tee $2
          br_if $block_4
          local.get $0
          i32.eqz
          br_if $block
          i32.const 0
          local.set $2
          local.get $0
          i32.load
          i32.const 0
          i32.const -1
          local.get $1
          i32.const 1
          i32.add
          local.tee $1
          i32.shl
          local.get $1
          i32.const 31
          i32.gt_u
          select
          i32.and
          local.tee $1
          i32.eqz
          br_if $block_3
          local.get $1
          call $34
          local.tee $1
          i32.const 2
          i32.shl
          local.get $3
          i32.add
          i32.load
          local.set $2
        end ;; $block_4
        local.get $0
        local.get $2
        call $34
        local.get $1
        i32.const 4
        i32.shl
        i32.add
        i32.const 2
        i32.shl
        i32.add
        i32.const 96
        i32.add
        i32.load
        local.set $2
      end ;; $block_3
      local.get $2
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $83 (type $6)
    (param $0 i32)
    (result i32)
    (local $1 i32)
    (local $2 i32)
    block $block
      local.get $0
      i32.const 16
      i32.shr_u
      local.get $0
      local.get $0
      i32.const 65535
      i32.gt_u
      local.tee $1
      select
      local.tee $0
      i32.const 8
      i32.shr_u
      local.get $0
      local.get $0
      i32.const 255
      i32.gt_u
      local.tee $2
      select
      local.tee $0
      i32.const 256
      i32.lt_u
      br_if $block
      call $68
      unreachable
    end ;; $block
    i32.const 32
    local.get $1
    i32.const 4
    i32.shl
    local.tee $1
    i32.const 8
    i32.or
    local.get $1
    local.get $2
    select
    local.get $0
    i32.const 65568
    i32.add
    i32.load8_u
    i32.add
    i32.sub
    )
  
  (func $84 (type $16)
    (param $0 i32)
    (param $1 i32)
    (param $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    block $block
      local.get $0
      i32.eqz
      br_if $block
      local.get $2
      i32.const -16
      i32.and
      local.set $2
      local.get $1
      i32.const 19
      i32.add
      i32.const -16
      i32.and
      local.tee $3
      i32.const -4
      i32.add
      local.set $1
      i32.const 1
      local.set $4
      block $block_0
        local.get $0
        i32.load
        i32.load offset=1568
        local.tee $5
        i32.eqz
        br_if $block_0
        local.get $3
        i32.const -20
        i32.add
        local.tee $3
        local.get $5
        i32.ne
        br_if $block_0
        local.get $5
        i32.load
        i32.const 2
        i32.and
        i32.const 1
        i32.or
        local.set $4
        local.get $3
        local.set $1
      end ;; $block_0
      block $block_1
        local.get $2
        local.get $1
        i32.sub
        local.tee $2
        i32.const 19
        i32.le_u
        br_if $block_1
        local.get $1
        i32.const 4
        i32.add
        local.tee $3
        i64.const 0
        i64.store align=4
        local.get $1
        local.get $4
        local.get $2
        i32.const -8
        i32.add
        local.tee $5
        i32.or
        i32.store
        local.get $3
        local.get $5
        i32.add
        local.tee $4
        i32.eqz
        br_if $block
        local.get $4
        i32.const 2
        i32.store
        local.get $0
        i32.load
        local.get $4
        i32.store offset=1568
        local.get $0
        i32.const 56
        i32.add
        local.tee $4
        local.get $4
        i64.load
        local.get $5
        i64.extend_i32_u
        i64.add
        i64.store
        local.get $0
        local.get $0
        i64.load offset=32
        local.get $2
        i64.extend_i32_u
        i64.add
        i64.store offset=32
        local.get $0
        local.get $1
        call $86
        return
      end ;; $block_1
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $85 (type $17)
    (param $0 i32)
    (param $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    block $block
      local.get $0
      i32.eqz
      br_if $block
      local.get $1
      i32.eqz
      br_if $block
      local.get $0
      i32.load
      local.set $2
      block $block_0
        block $block_1
          local.get $1
          i32.load
          local.tee $3
          i32.const -4
          i32.and
          local.tee $0
          i32.const 256
          i32.ge_u
          br_if $block_1
          local.get $3
          i32.const 4
          i32.shr_u
          local.set $3
          i32.const 0
          local.set $4
          br $block_0
        end ;; $block_1
        local.get $0
        i32.const 1073741820
        local.get $0
        i32.const 1073741820
        i32.lt_u
        select
        local.set $0
        i32.const 16
        local.get $0
        i32.const 27
        local.get $0
        call $83
        local.tee $5
        i32.sub
        local.tee $3
        i32.shr_u
        i32.const 16
        i32.xor
        local.get $3
        i32.const 31
        i32.gt_u
        select
        local.set $3
        i32.const 24
        local.get $5
        i32.sub
        local.set $4
      end ;; $block_0
      local.get $1
      i32.load offset=8
      local.set $0
      block $block_2
        local.get $1
        i32.load offset=4
        local.tee $5
        i32.eqz
        br_if $block_2
        local.get $5
        local.get $0
        i32.store offset=8
      end ;; $block_2
      block $block_3
        local.get $0
        i32.eqz
        br_if $block_3
        local.get $0
        local.get $5
        i32.store offset=4
      end ;; $block_3
      block $block_4
        local.get $2
        local.get $4
        i32.const 4
        i32.shl
        local.get $3
        i32.add
        i32.const 2
        i32.shl
        i32.add
        i32.const 96
        i32.add
        local.tee $5
        i32.load
        local.get $1
        i32.ne
        br_if $block_4
        local.get $5
        local.get $0
        i32.store
        local.get $0
        br_if $block_4
        local.get $2
        local.get $4
        i32.const 2
        i32.shl
        i32.add
        i32.const 4
        i32.add
        local.tee $1
        local.get $1
        i32.load
        i32.const -1
        i32.const -2
        local.get $3
        i32.rotl
        local.get $3
        i32.const 31
        i32.gt_u
        select
        i32.and
        local.tee $1
        i32.store
        local.get $1
        br_if $block_4
        local.get $2
        i32.eqz
        br_if $block
        local.get $2
        local.get $2
        i32.load
        i32.const -1
        i32.const -2
        local.get $4
        i32.rotl
        local.get $4
        i32.const 31
        i32.gt_u
        select
        i32.and
        i32.store
      end ;; $block_4
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $86 (type $17)
    (param $0 i32)
    (param $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    block $block
      local.get $0
      i32.eqz
      br_if $block
      local.get $1
      i32.eqz
      br_if $block
      local.get $1
      i32.load
      local.tee $2
      i32.const -4
      i32.and
      local.get $1
      i32.const 4
      i32.add
      local.tee $3
      i32.add
      local.tee $4
      i32.eqz
      br_if $block
      local.get $0
      i32.load
      local.set $5
      block $block_0
        local.get $4
        i32.load
        local.tee $6
        i32.const 1
        i32.and
        i32.eqz
        br_if $block_0
        local.get $0
        local.get $4
        call $85
        local.get $1
        local.get $2
        local.get $6
        i32.const -4
        i32.and
        i32.add
        i32.const 4
        i32.add
        local.tee $2
        i32.store
        local.get $2
        i32.const -4
        i32.and
        local.get $3
        i32.add
        local.tee $4
        i32.eqz
        br_if $block
        local.get $4
        i32.load
        local.set $6
      end ;; $block_0
      block $block_1
        local.get $2
        i32.const 2
        i32.and
        i32.eqz
        br_if $block_1
        local.get $1
        i32.const -4
        i32.add
        i32.load
        local.tee $1
        i32.eqz
        br_if $block
        local.get $1
        i32.load
        local.set $3
        local.get $0
        local.get $1
        call $85
        local.get $1
        local.get $3
        local.get $2
        i32.const 4
        i32.add
        i32.const -4
        i32.and
        i32.add
        local.tee $2
        i32.store
      end ;; $block_1
      local.get $4
      i32.const -4
      i32.add
      local.get $1
      i32.store
      local.get $4
      local.get $6
      i32.const 2
      i32.or
      i32.store
      block $block_2
        block $block_3
          local.get $2
          i32.const -4
          i32.and
          local.tee $4
          i32.const 256
          i32.ge_u
          br_if $block_3
          local.get $2
          i32.const 4
          i32.shr_u
          local.set $4
          i32.const 0
          local.set $2
          br $block_2
        end ;; $block_3
        local.get $4
        i32.const 1073741820
        local.get $4
        i32.const 1073741820
        i32.lt_u
        select
        local.set $2
        i32.const 16
        local.get $2
        i32.const 27
        local.get $2
        call $83
        local.tee $0
        i32.sub
        local.tee $4
        i32.shr_u
        i32.const 16
        i32.xor
        local.get $4
        i32.const 31
        i32.gt_u
        select
        local.set $4
        i32.const 24
        local.get $0
        i32.sub
        local.set $2
      end ;; $block_2
      local.get $1
      local.get $5
      local.get $2
      i32.const 4
      i32.shl
      local.get $4
      i32.add
      i32.const 2
      i32.shl
      i32.add
      i32.const 96
      i32.add
      local.tee $6
      i32.load
      local.tee $0
      i32.store offset=8
      local.get $1
      i32.const 0
      i32.store offset=4
      block $block_4
        local.get $0
        i32.eqz
        br_if $block_4
        local.get $0
        local.get $1
        i32.store offset=4
      end ;; $block_4
      local.get $6
      local.get $1
      i32.store
      local.get $5
      i32.eqz
      br_if $block
      local.get $5
      local.get $5
      i32.load
      i32.const 0
      i32.const 1
      local.get $2
      i32.shl
      local.get $2
      i32.const 31
      i32.gt_u
      select
      i32.or
      i32.store
      local.get $5
      local.get $2
      i32.const 2
      i32.shl
      i32.add
      i32.const 4
      i32.add
      local.tee $1
      local.get $1
      i32.load
      i32.const 0
      i32.const 1
      local.get $4
      i32.shl
      local.get $4
      i32.const 31
      i32.gt_u
      select
      i32.or
      i32.store
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $87 (type $5)
    i32.const 66072
    i32.const 23
    call $100
    )
  
  (func $88 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    block $block
      local.get $0
      local.get $1
      call $80
      local.tee $0
      br_if $block
      i32.const 0
      return
    end ;; $block
    local.get $0
    i32.const 4
    i32.add
    i32.const 0
    local.get $1
    call $145
    )
  
  (func $89 (type $6)
    (param $0 i32)
    (result i32)
    block $block
      block $block_0
        block $block_1
          local.get $0
          i32.eqz
          br_if $block_1
          local.get $0
          i32.const 15
          i32.and
          br_if $block_1
          local.get $0
          i32.const -4
          i32.add
          local.tee $0
          i32.eqz
          br_if $block_0
          local.get $0
          i32.load8_u
          i32.const 1
          i32.and
          i32.eqz
          br_if $block
        end ;; $block_1
        i32.const 66000
        call $81
        unreachable
      end ;; $block_0
      call $87
      unreachable
    end ;; $block
    local.get $0
    )
  
  (func $90 (type $17)
    (param $0 i32)
    (param $1 i32)
    (local $2 i32)
    (local $3 i64)
    block $block
      local.get $1
      i32.eqz
      br_if $block
      local.get $0
      i32.eqz
      br_if $block
      local.get $0
      i32.const 56
      i32.add
      local.tee $2
      local.get $2
      i64.load
      local.get $1
      i64.load32_u
      i64.const 4294967292
      i64.and
      local.tee $3
      i64.add
      i64.store
      local.get $0
      i32.const 40
      i32.add
      local.tee $2
      local.get $2
      i64.load
      local.get $3
      i64.sub
      i64.store
      local.get $0
      i32.const 64
      i32.add
      local.tee $2
      local.get $2
      i32.load
      i32.const -1
      i32.add
      i32.store
      local.get $1
      local.get $1
      i32.load
      i32.const 1
      i32.or
      i32.store
      local.get $0
      local.get $1
      call $86
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $91 (type $6)
    (param $0 i32)
    (result i32)
    block $block
      local.get $0
      i32.const -4
      i32.add
      local.tee $0
      br_if $block
      call $87
      unreachable
    end ;; $block
    local.get $0
    i32.load
    i32.const -4
    i32.and
    )
  
  (func $92 (type $6)
    (param $0 i32)
    (result i32)
    local.get $0
    i32.const 255
    i32.and
    i32.const -2128831035
    i32.xor
    i32.const 16777619
    i32.mul
    local.get $0
    i32.const 8
    i32.shr_u
    i32.const 255
    i32.and
    i32.xor
    i32.const 16777619
    i32.mul
    local.get $0
    i32.const 16
    i32.shr_u
    i32.const 255
    i32.and
    i32.xor
    i32.const 16777619
    i32.mul
    local.get $0
    i32.const 24
    i32.shr_u
    i32.xor
    i32.const 16777619
    i32.mul
    )
  
  (func $93 (type $18)
    (result i32)
    (local $0 i32)
    block $block
      i32.const 0
      i32.load offset=66824
      local.tee $0
      i32.eqz
      br_if $block
      i32.const 0
      local.get $0
      i32.load
      i32.store offset=66824
      block $block_0
        i32.const 0
        i32.load offset=66828
        local.get $0
        i32.ne
        br_if $block_0
        i32.const 0
        i32.const 0
        i32.store offset=66828
      end ;; $block_0
      local.get $0
      i32.const 0
      i32.store
      local.get $0
      return
    end ;; $block
    i32.const 0
    )
  
  (func $94 (type $17)
    (param $0 i32)
    (param $1 i32)
    (local $2 i32)
    block $block
      local.get $0
      i32.eqz
      br_if $block
      block $block_0
        local.get $0
        i32.load offset=4
        local.tee $2
        i32.eqz
        br_if $block_0
        local.get $2
        local.get $1
        i32.store
      end ;; $block_0
      local.get $0
      local.get $1
      i32.store offset=4
      local.get $1
      i32.eqz
      br_if $block
      local.get $1
      i32.const 0
      i32.store
      block $block_1
        local.get $0
        i32.load
        br_if $block_1
        local.get $0
        local.get $1
        i32.store
      end ;; $block_1
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $95 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    (local $2 i32)
    block $block
      local.get $0
      br_if $block
      call $87
      unreachable
    end ;; $block
    local.get $0
    i32.load offset=16
    local.set $2
    local.get $0
    local.get $1
    i32.store offset=16
    local.get $2
    )
  
  (func $96 (type $18)
    (result i32)
    (local $0 i32)
    i32.const 24
    call $97
    local.tee $0
    i32.const 66488
    i32.store offset=16
    local.get $0
    )
  
  (func $97 (type $6)
    (param $0 i32)
    (result i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i64)
    block $block
      block $block_0
        local.get $0
        i32.const 1073741813
        i32.ge_u
        br_if $block_0
        i32.const 0
        local.set $1
        i32.const 0
        i32.load offset=66860
        local.set $2
        i32.const 0
        i32.load offset=66856
        local.get $0
        i32.const 12
        i32.add
        call $88
        local.tee $3
        call $91
        local.set $4
        block $block_1
          local.get $3
          i32.eqz
          br_if $block_1
          local.get $3
          local.get $0
          i32.store offset=8
          local.get $3
          local.get $4
          i32.store
          local.get $3
          i32.const 0
          i32.store offset=4
          local.get $2
          i32.eqz
          br_if $block
          local.get $2
          i32.const 344
          i32.add
          local.tee $0
          local.get $0
          i32.load
          local.get $4
          i32.add
          i32.store
          local.get $3
          i64.load32_u
          local.set $5
          local.get $2
          i32.const 72
          i32.add
          local.tee $0
          local.get $0
          i64.load
          i64.const 1
          i64.add
          i64.store
          local.get $2
          i32.const 80
          i32.add
          local.tee $0
          local.get $0
          i64.load
          i64.const 1
          i64.add
          i64.store
          local.get $2
          i32.const 88
          i32.add
          local.tee $0
          local.get $5
          local.get $0
          i64.load
          i64.add
          i64.store
          local.get $2
          local.get $3
          i32.const 0
          call $98
          drop
          block $block_2
            local.get $3
            local.get $2
            i32.load offset=28
            i32.ge_u
            br_if $block_2
            local.get $2
            local.get $3
            i32.store offset=28
          end ;; $block_2
          block $block_3
            local.get $3
            local.get $2
            i32.load offset=32
            i32.le_u
            br_if $block_3
            local.get $2
            local.get $3
            i32.store offset=32
          end ;; $block_3
          local.get $3
          i32.const 12
          i32.add
          local.set $1
        end ;; $block_1
        local.get $1
        return
      end ;; $block_0
      i32.const 66216
      call $81
      unreachable
    end ;; $block
    call $87
    unreachable
    )
  
  (func $98 (type $19)
    (param $0 i32)
    (param $1 i32)
    (param $2 i32)
    (result i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    (local $9 f32)
    global.get $24
    i32.const 32
    i32.sub
    local.tee $3
    global.set $24
    block $block
      local.get $0
      i32.eqz
      br_if $block
      i32.const 0
      local.set $4
      local.get $1
      call $92
      local.get $0
      i32.load offset=8
      i32.rem_u
      i32.const 3
      i32.shl
      local.get $0
      i32.load
      i32.add
      local.tee $5
      local.set $6
      block $block_0
        block $block_1
          loop $loop
            local.get $6
            i32.eqz
            br_if $block
            local.get $6
            i32.load
            local.tee $7
            i32.eqz
            br_if $block_0
            block $block_2
              local.get $7
              local.get $1
              i32.eq
              br_if $block_2
              block $block_3
                block $block_4
                  local.get $6
                  i32.load offset=4
                  local.tee $8
                  local.get $4
                  i32.ge_s
                  br_if $block_4
                  local.get $6
                  local.get $4
                  i32.store offset=4
                  local.get $6
                  local.get $1
                  i32.store
                  local.get $7
                  local.set $1
                  br $block_3
                end ;; $block_4
                local.get $4
                local.set $8
              end ;; $block_3
              block $block_5
                local.get $6
                i32.const 8
                i32.add
                local.tee $6
                local.get $0
                i32.load offset=4
                i32.lt_u
                br_if $block_5
                local.get $0
                i32.load
                local.set $6
              end ;; $block_5
              local.get $6
              local.get $5
              i32.eq
              br_if $block_1
              local.get $8
              i32.const 1
              i32.add
              local.tee $4
              local.get $0
              i32.load offset=16
              i32.gt_s
              br_if $block_1
              br $loop
            end ;; $block_2
          end ;; $loop
          local.get $6
          local.get $4
          i32.store offset=4
          local.get $3
          i32.const 32
          i32.add
          global.set $24
          i32.const 1
          return
        end ;; $block_1
        block $block_6
          local.get $2
          i32.const 6
          i32.ge_s
          br_if $block_6
          block $block_7
            local.get $0
            f32.load offset=24
            local.tee $9
            f32.const 0x1.000000p-0
            f32.le
            i32.const 1
            i32.xor
            br_if $block_7
            local.get $0
            i32.const 1073741824
            i32.store offset=24
            f32.const 0x1.000000p+1
            local.set $9
          end ;; $block_7
          i32.const -1
          i32.const 0
          local.get $9
          local.get $0
          i32.load offset=8
          f32.convert_i32_u
          f32.mul
          local.tee $9
          f32.const 0x0.000000p-127
          f32.ge
          local.tee $6
          select
          local.set $4
          local.get $9
          f32.const 0x1.fffffcp+31
          f32.le
          local.set $7
          block $block_8
            block $block_9
              local.get $9
              f32.const 0x1.000000p+32
              f32.lt
              local.get $6
              i32.and
              i32.eqz
              br_if $block_9
              local.get $9
              i32.trunc_f32_u
              local.set $8
              br $block_8
            end ;; $block_9
            i32.const 0
            local.set $8
          end ;; $block_8
          local.get $8
          local.get $4
          local.get $7
          select
          local.get $4
          local.get $6
          select
          local.tee $4
          i32.const 3
          i32.shl
          local.tee $7
          call $118
          local.set $6
          local.get $3
          i32.const 0
          i32.store offset=24
          local.get $3
          local.get $4
          i32.store offset=12
          local.get $3
          local.get $6
          i32.store offset=4
          local.get $3
          i32.const 0
          i32.store offset=16
          local.get $3
          local.get $0
          i32.load offset=24
          i32.store offset=28
          local.get $3
          local.get $7
          local.get $6
          i32.add
          i32.store offset=8
          local.get $3
          local.get $0
          i32.load offset=16
          i32.store offset=20
          local.get $0
          i32.load
          local.set $6
          block $block_10
            loop $loop_0
              local.get $6
              local.get $0
              i32.load offset=4
              i32.ge_u
              br_if $block_10
              local.get $6
              i32.load
              local.set $4
              local.get $6
              i32.const 8
              i32.add
              local.set $6
              local.get $3
              i32.const 4
              i32.add
              local.get $4
              i32.const 0
              call $98
              i32.const 1
              i32.and
              i32.eqz
              br_if $block_6
              br $loop_0
            end ;; $loop_0
          end ;; $block_10
          local.get $0
          i32.load
          call $108
          local.get $0
          local.get $3
          i32.load offset=28
          i32.store offset=24
          local.get $0
          local.get $3
          i64.load offset=20 align=4
          i64.store offset=16 align=4
          local.get $0
          local.get $3
          i64.load offset=12 align=4
          i64.store offset=8 align=4
          local.get $0
          local.get $3
          i64.load offset=4 align=4
          i64.store align=4
          local.get $0
          local.get $1
          local.get $2
          i32.const 1
          i32.add
          call $98
          local.set $6
          local.get $3
          i32.const 32
          i32.add
          global.set $24
          local.get $6
          return
        end ;; $block_6
        local.get $3
        i32.const 32
        i32.add
        global.set $24
        i32.const 0
        return
      end ;; $block_0
      local.get $6
      local.get $4
      i32.store offset=4
      local.get $6
      local.get $1
      i32.store
      local.get $0
      local.get $0
      i32.load offset=12
      i32.const 1
      i32.add
      i32.store offset=12
      local.get $3
      i32.const 32
      i32.add
      global.set $24
      i32.const 1
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $99 (type $5)
    i32.const 66135
    i32.const 18
    call $100
    )
  
  (func $100 (type $17)
    (param $0 i32)
    (param $1 i32)
    i32.const 66095
    i32.const 22
    call $101
    local.get $0
    local.get $1
    call $101
    call $102
    unreachable
    unreachable
    )
  
  (func $101 (type $17)
    (param $0 i32)
    (param $1 i32)
    local.get $1
    i32.const 0
    local.get $1
    i32.const 0
    i32.gt_s
    select
    local.set $1
    block $block
      loop $loop
        local.get $1
        i32.eqz
        br_if $block
        local.get $0
        i32.load8_u
        call $103
        local.get $0
        i32.const 1
        i32.add
        local.set $0
        local.get $1
        i32.const -1
        i32.add
        local.set $1
        br $loop
      end ;; $loop
    end ;; $block
    )
  
  (func $102 (type $5)
    i32.const 13
    call $103
    i32.const 10
    call $103
    )
  
  (func $103 (type $0)
    (param $0 i32)
    (local $1 i32)
    (local $2 i32)
    block $block
      i32.const 0
      i32.load offset=66564
      local.tee $1
      i32.const 255
      i32.le_u
      br_if $block
      call $68
      unreachable
    end ;; $block
    i32.const 0
    local.get $1
    i32.const 1
    i32.add
    local.tee $2
    i32.store offset=66564
    local.get $1
    i32.const 66568
    i32.add
    local.get $0
    i32.store8
    block $block_0
      block $block_1
        local.get $0
        i32.const 255
        i32.and
        i32.const 10
        i32.eq
        br_if $block_1
        local.get $1
        i32.const 255
        i32.ne
        br_if $block_0
      end ;; $block_1
      i32.const 0
      local.get $2
      i32.store offset=66520
      i32.const 1
      i32.const 66516
      i32.const 1
      i32.const 66836
      call $28
      drop
      i32.const 0
      i32.const 0
      i32.store offset=66564
    end ;; $block_0
    )
  
  (func $104 (type $0)
    (param $0 i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    i32.const 66462
    i32.const 6
    call $101
    i32.const 32
    call $103
    local.get $0
    i64.extend_i32_u
    call $105
    call $102
    block $block
      i32.const 0
      i32.load offset=66860
      local.tee $1
      i32.eqz
      br_if $block
      block $block_0
        local.get $1
        local.get $0
        i32.const -12
        i32.add
        local.tee $0
        call $106
        i32.const 1
        i32.and
        i32.eqz
        br_if $block_0
        local.get $0
        call $107
        local.set $2
        local.get $1
        i32.const 344
        i32.add
        local.tee $3
        local.get $3
        i32.load
        local.get $2
        i32.sub
        i32.store
        local.get $1
        i32.const 104
        i32.add
        local.tee $3
        local.get $3
        i64.load
        local.get $2
        i64.extend_i32_u
        i64.add
        i64.store
        local.get $1
        i32.const 72
        i32.add
        local.tee $2
        local.get $2
        i64.load
        i64.const -1
        i64.add
        i64.store
        local.get $1
        i32.const 96
        i32.add
        local.tee $1
        local.get $1
        i64.load
        i64.const 1
        i64.add
        i64.store
        local.get $0
        call $108
      end ;; $block_0
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $105 (type $4)
    (param $0 i64)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i64)
    (local $5 i32)
    global.get $24
    i32.const 32
    i32.sub
    local.tee $1
    global.set $24
    local.get $1
    i32.const 24
    i32.add
    i32.const 0
    i32.store
    local.get $1
    i32.const 16
    i32.add
    i64.const 0
    i64.store
    local.get $1
    i64.const 0
    i64.store offset=8
    i32.const 19
    local.set $2
    i32.const 19
    local.set $3
    loop $loop
      block $block
        local.get $3
        i32.const -1
        i32.gt_s
        br_if $block
        local.get $2
        i32.const 20
        local.get $2
        i32.const 20
        i32.gt_s
        select
        local.get $2
        i32.sub
        local.set $3
        local.get $1
        i32.const 8
        i32.add
        local.get $2
        i32.add
        local.set $2
        block $block_0
          loop $loop_0
            local.get $3
            i32.eqz
            br_if $block_0
            local.get $2
            i32.load8_u
            call $103
            local.get $2
            i32.const 1
            i32.add
            local.set $2
            local.get $3
            i32.const -1
            i32.add
            local.set $3
            br $loop_0
          end ;; $loop_0
        end ;; $block_0
        local.get $1
        i32.const 32
        i32.add
        global.set $24
        return
      end ;; $block
      local.get $1
      i32.const 8
      i32.add
      local.get $3
      i32.add
      local.get $0
      local.get $0
      i64.const 10
      i64.div_u
      local.tee $4
      i64.const 10
      i64.mul
      i64.sub
      i32.wrap_i64
      i32.const 48
      i32.or
      local.tee $5
      i32.store8
      local.get $2
      local.get $3
      local.get $5
      i32.const 255
      i32.and
      i32.const 48
      i32.eq
      select
      local.set $2
      local.get $3
      i32.const -1
      i32.add
      local.set $3
      local.get $4
      local.set $0
      br $loop
    end ;; $loop
    )
  
  (func $106 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    block $block
      block $block_0
        block $block_1
          local.get $1
          i32.eqz
          br_if $block_1
          local.get $0
          i32.eqz
          br_if $block
          local.get $1
          call $92
          local.get $0
          i32.load offset=8
          i32.rem_u
          i32.const 3
          i32.shl
          local.get $0
          i32.load
          local.tee $2
          i32.add
          local.tee $3
          local.set $4
          loop $loop
            local.get $4
            i32.eqz
            br_if $block
            local.get $4
            i32.load
            local.tee $5
            i32.eqz
            br_if $block_1
            block $block_2
              local.get $5
              local.get $1
              i32.ne
              br_if $block_2
              br $block_0
            end ;; $block_2
            local.get $4
            i32.const 8
            i32.add
            local.tee $4
            local.get $2
            local.get $4
            local.get $0
            i32.load offset=4
            i32.lt_u
            select
            local.tee $4
            local.get $3
            i32.ne
            br_if $loop
          end ;; $loop
        end ;; $block_1
        i32.const 0
        return
      end ;; $block_0
      block $block_3
        loop $loop_0
          block $block_4
            local.get $4
            i32.const 8
            i32.add
            local.tee $5
            local.get $0
            i32.load offset=4
            i32.lt_u
            br_if $block_4
            local.get $0
            i32.load
            local.tee $5
            local.get $4
            i32.eq
            br_if $block_3
          end ;; $block_4
          local.get $5
          i32.eqz
          br_if $block
          local.get $5
          i32.load
          i32.eqz
          br_if $block_3
          local.get $5
          i32.load offset=4
          local.tee $1
          i32.const 1
          i32.lt_s
          br_if $block_3
          local.get $5
          local.get $1
          i32.const -1
          i32.add
          i32.store offset=4
          local.get $4
          local.get $5
          i64.load align=4
          i64.store align=4
          local.get $5
          local.set $4
          br $loop_0
        end ;; $loop_0
      end ;; $block_3
      local.get $4
      i64.const 0
      i64.store align=4
      local.get $0
      local.get $0
      i32.load offset=12
      i32.const -1
      i32.add
      i32.store offset=12
      i32.const 1
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $107 (type $6)
    (param $0 i32)
    (result i32)
    block $block
      local.get $0
      br_if $block
      call $87
      unreachable
    end ;; $block
    local.get $0
    i32.load
    )
  
  (func $108 (type $0)
    (param $0 i32)
    i32.const 0
    i32.load offset=66856
    local.get $0
    call $89
    call $90
    )
  
  (func $109 (type $6)
    (param $0 i32)
    (result i32)
    (local $1 i32)
    (local $2 i32)
    block $block
      i32.const 0
      i32.load offset=66548
      br_if $block
      i32.const 0
      return
    end ;; $block
    local.get $0
    i32.const 0
    i32.load offset=66552
    i32.const 0
    i32.load offset=66544
    i32.const 0
    i32.load offset=66524
    local.tee $1
    i32.mul
    i32.add
    local.tee $2
    local.get $1
    call $144
    drop
    local.get $2
    i32.const 0
    i32.const 0
    i32.load offset=66524
    call $145
    drop
    i32.const 0
    i32.const 0
    i32.load offset=66548
    i32.const -1
    i32.add
    i32.store offset=66548
    i32.const 0
    i32.const 0
    i32.const 0
    i32.load offset=66544
    i32.const 1
    i32.add
    local.tee $0
    local.get $0
    i32.const 0
    i32.load offset=66528
    i32.eq
    select
    i32.store offset=66544
    i32.const 1
    )
  
  (func $110 (type $18)
    (result i32)
    (local $0 i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    (local $9 i32)
    (local $10 i32)
    (local $11 i32)
    (local $12 i32)
    block $block
      block $block_0
        i32.const 0
        i32.load offset=66536
        local.tee $0
        i32.eqz
        br_if $block_0
        i32.const 0
        local.get $0
        i32.load
        i32.store offset=66536
        local.get $0
        i32.load offset=4
        local.tee $1
        i32.eqz
        br_if $block_0
        block $block_1
          local.get $0
          i32.load offset=8
          local.tee $2
          br_if $block_1
          local.get $1
          i32.load offset=4
          local.set $3
          br $block
        end ;; $block_1
        local.get $2
        i32.load offset=4
        local.set $3
        local.get $1
        local.get $2
        i32.store offset=4
        local.get $0
        i32.load offset=12
        local.tee $4
        i32.eqz
        br_if $block
        i32.const -1
        local.set $5
        local.get $0
        i32.const 16
        i32.add
        i32.load
        local.tee $2
        i32.const 0
        local.get $2
        i32.const 0
        i32.gt_s
        select
        i32.const -1
        i32.add
        local.set $6
        local.get $2
        i32.const -1
        i32.add
        local.set $7
        loop $loop
          local.get $5
          local.get $6
          i32.eq
          br_if $block
          block $block_2
            block $block_3
              local.get $5
              local.get $7
              i32.eq
              br_if $block_3
              local.get $5
              i32.const 1
              i32.add
              local.tee $5
              local.get $0
              i32.load offset=16
              i32.ge_u
              br_if $block_3
              local.get $0
              i32.load offset=12
              local.tee $8
              local.get $5
              i32.const 24
              i32.mul
              local.tee $2
              i32.add
              local.tee $9
              local.get $0
              i32.eq
              br_if $loop
              local.get $4
              local.get $2
              i32.add
              local.tee $2
              i32.load offset=4
              i32.eqz
              br_if $loop
              local.get $2
              i32.load offset=8
              local.tee $10
              i32.eqz
              br_if $block_0
              local.get $10
              i32.load
              local.tee $11
              i32.eqz
              br_if $block_0
              local.get $11
              i32.load offset=12
              local.tee $12
              local.set $2
              block $block_4
                block $block_5
                  local.get $12
                  local.get $9
                  i32.ne
                  br_if $block_5
                  local.get $8
                  i32.eqz
                  br_if $block_0
                  local.get $9
                  i32.load
                  local.set $12
                  br $block_4
                end ;; $block_5
                loop $loop_0
                  local.get $2
                  local.tee $1
                  i32.eqz
                  br_if $block_4
                  local.get $1
                  i32.load
                  local.tee $2
                  local.get $9
                  i32.ne
                  br_if $loop_0
                end ;; $loop_0
                local.get $8
                i32.eqz
                br_if $block_0
                local.get $1
                local.get $9
                i32.load
                i32.store
              end ;; $block_4
              local.get $11
              local.get $12
              i32.store offset=12
              local.get $10
              i32.load
              local.tee $2
              i32.eqz
              br_if $block_0
              local.get $2
              i32.load offset=12
              br_if $loop
              block $block_6
                local.get $10
                i32.load offset=4
                br_if $block_6
                local.get $2
                i32.load8_u offset=8
                i32.const 4
                i32.eq
                br_if $loop
                br $block_2
              end ;; $block_6
              local.get $2
              i32.load offset=24
              i32.eqz
              br_if $block_2
              local.get $2
              i32.const 3
              i32.store8 offset=8
              br $loop
            end ;; $block_3
            call $68
            unreachable
          end ;; $block_2
          local.get $2
          i32.const 0
          i32.store8 offset=8
          br $loop
        end ;; $loop
      end ;; $block_0
      call $87
      unreachable
    end ;; $block
    i32.const 66824
    local.get $0
    i32.load offset=4
    call $94
    local.get $3
    )
  
  (func $111 (type $0)
    (param $0 i32)
    local.get $0
    i32.load
    local.get $0
    i32.load offset=4
    call $101
    )
  
  (func $112 (type $4)
    (param $0 i64)
    block $block
      local.get $0
      i64.const -1
      i64.gt_s
      br_if $block
      i32.const 45
      call $103
      i64.const 0
      local.get $0
      i64.sub
      local.set $0
    end ;; $block
    local.get $0
    call $105
    )
  
  (func $113 (type $20)
    (param $0 i32)
    (param $1 i32)
    (param $2 i32)
    (param $3 i32)
    (param $4 i32)
    (result i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    (local $9 i32)
    (local $10 i32)
    i32.const 24
    call $97
    local.tee $5
    local.get $4
    i32.store offset=8
    local.get $5
    i32.const 1
    i32.store offset=4
    local.get $5
    i32.const 2
    i32.store
    local.get $5
    local.get $4
    local.get $5
    call $95
    local.tee $6
    i32.store offset=12
    block $block
      local.get $4
      i32.eqz
      br_if $block
      local.get $5
      local.get $4
      i32.load offset=4
      local.tee $7
      i32.store offset=16
      i32.const 0
      local.set $8
      i32.const 0
      local.set $9
      block $block_0
        block $block_1
          block $block_2
            block $block_3
              block $block_4
                block $block_5
                  block $block_6
                    block $block_7
                      i32.const 0
                      i32.load8_u offset=66532
                      br_table
                        $block_0 $block_0 $block_6 $block_6 $block_3
                        $block_7 ;; default
                    end ;; $block_7
                    i32.const 66153
                    i32.const 21
                    call $100
                    i32.const 66174
                    i32.const 11
                    call $100
                    i32.const 0
                    local.set $8
                    br $block_5
                  end ;; $block_6
                  local.get $1
                  call $109
                  local.set $10
                  i32.const 0
                  local.set $8
                  i32.const 0
                  i32.load offset=66536
                  local.set $9
                  local.get $10
                  i32.const 1
                  i32.and
                  br_if $block_2
                  local.get $9
                  br_if $block_4
                end ;; $block_5
                i32.const 0
                local.set $9
                br $block_0
              end ;; $block_4
              local.get $1
              call $110
              i32.const 0
              i32.load offset=66524
              call $144
              drop
              i32.const 1
              local.set $8
              i32.const 1
              local.set $9
              i32.const 0
              i32.load offset=66536
              br_if $block_0
              i32.const 0
              i32.const 0
              i32.store8 offset=66532
              br $block_1
            end ;; $block_3
            i32.const 1
            local.set $8
            i32.const 1
            local.set $9
            local.get $1
            call $109
            i32.const 1
            i32.and
            br_if $block_0
            i32.const 0
            local.set $9
            local.get $1
            i32.const 0
            i32.const 0
            i32.load offset=66524
            call $145
            drop
            i32.const 1
            local.set $8
            br $block_0
          end ;; $block_2
          block $block_8
            block $block_9
              local.get $9
              br_if $block_9
              i32.const 0
              i32.load offset=66548
              local.set $10
              br $block_8
            end ;; $block_9
            call $110
            local.set $9
            i32.const 0
            i32.load offset=66548
            local.set $10
            block $block_10
              i32.const 0
              i32.load offset=66528
              local.tee $8
              i32.eqz
              br_if $block_10
              local.get $10
              local.get $8
              i32.eq
              br_if $block_10
              i32.const 0
              i32.load offset=66552
              i32.const 0
              i32.load offset=66540
              i32.const 0
              i32.load offset=66524
              local.tee $8
              i32.mul
              i32.add
              local.get $9
              local.get $8
              call $144
              drop
              i32.const 0
              i32.const 0
              i32.load offset=66548
              i32.const 1
              i32.add
              local.tee $10
              i32.store offset=66548
              i32.const 0
              i32.const 0
              i32.const 0
              i32.load offset=66540
              i32.const 1
              i32.add
              local.tee $8
              local.get $8
              i32.const 0
              i32.load offset=66528
              i32.eq
              select
              i32.store offset=66540
            end ;; $block_10
            i32.const 0
            i32.load offset=66536
            br_if $block_8
            i32.const 0
            i32.const 3
            i32.store8 offset=66532
          end ;; $block_8
          i32.const 1
          local.set $8
          i32.const 1
          local.set $9
          local.get $10
          br_if $block_0
          i32.const 0
          i32.const 0
          i32.store8 offset=66532
        end ;; $block_1
        i32.const 1
        local.set $8
        i32.const 1
        local.set $9
      end ;; $block_0
      block $block_11
        local.get $8
        i32.const 1
        i32.and
        i32.eqz
        br_if $block_11
        local.get $7
        local.get $9
        i32.const 1
        i32.and
        i32.store8
        local.get $4
        local.get $6
        i32.store offset=16
        i32.const 66824
        local.get $4
        call $94
        local.get $5
        call $104
        local.get $4
        return
      end ;; $block_11
      i32.const 0
      i32.const 1
      i32.store8 offset=66532
      local.get $4
      i64.const 1
      i64.store offset=8
      local.get $4
      local.get $1
      i32.store offset=4
      local.get $2
      i32.eqz
      br_if $block
      i32.const 0
      i32.load offset=66536
      local.set $8
      local.get $2
      local.get $4
      i32.store offset=4
      local.get $2
      local.get $8
      i32.store
      local.get $2
      i32.const 16
      i32.add
      i64.const 0
      i64.store align=4
      local.get $2
      i64.const 0
      i64.store offset=8 align=4
      local.get $5
      i32.const 0
      i32.store8 offset=20
      i32.const 0
      local.get $2
      i32.store offset=66536
      local.get $4
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $114 (type $0)
    (param $0 i32)
    local.get $0
    call $104
    )
  
  (func $115 (type $0)
    (param $0 i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i64)
    local.get $0
    i32.load offset=12
    local.set $1
    local.get $0
    i32.load offset=16
    local.set $2
    local.get $0
    i32.load offset=8
    local.tee $3
    i32.const 0
    i32.store offset=4
    local.get $3
    i64.load offset=8
    local.set $4
    local.get $3
    i64.const 0
    i64.store offset=8
    local.get $2
    local.get $4
    i64.const 1
    i64.eq
    i32.store8
    local.get $3
    local.get $1
    i32.store offset=16
    i32.const 66824
    local.get $3
    call $94
    local.get $0
    call $104
    )
  
  (func $116 (type $21)
    (result i64)
    (local $0 i32)
    (local $1 i64)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $0
    global.set $24
    local.get $0
    i64.const 0
    i64.store offset=8
    i32.const 0
    i64.const 1000
    local.get $0
    i32.const 8
    i32.add
    call $29
    drop
    local.get $0
    i64.load offset=8
    local.set $1
    local.get $0
    i32.const 16
    i32.add
    global.set $24
    local.get $1
    )
  
  (func $117 (type $5)
    (local $0 i32)
    (local $1 i64)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $0
    global.set $24
    i32.const 0
    memory.size
    i32.const 16
    i32.shl
    i32.store offset=66560
    call $33
    i32.const 0
    i64.load32_u offset=66560
    local.set $1
    i32.const 66452
    i32.const 10
    call $101
    i32.const 32
    call $103
    i32.const 66876
    i64.extend_i32_u
    call $105
    i32.const 32
    call $103
    local.get $1
    call $105
    call $102
    block $block
      block $block_0
        i32.const 0
        i32.load offset=66856
        br_if $block_0
        block $block_1
          memory.size
          i32.const 16
          i32.shl
          local.tee $2
          i32.const 68567
          i32.ge_u
          br_if $block_1
          i32.const 1
          memory.grow
          drop
          memory.size
          i32.const 16
          i32.shl
          local.set $2
        end ;; $block_1
        i32.const 66879
        i32.const -4
        i32.and
        i32.const 0
        i32.const 96
        call $145
        local.tee $3
        i32.const 80
        i32.add
        i32.const 1
        i32.store
        local.get $3
        i32.const 68
        i32.add
        i32.const 1
        i32.store
        local.get $3
        i32.const 8
        i32.add
        local.get $2
        i32.store
        local.get $3
        i32.const 4
        i32.add
        local.get $3
        i32.store
        local.get $3
        i32.const 20
        i32.add
        i32.const 1
        i32.store
        local.get $3
        local.get $3
        i32.const 15
        i32.add
        i32.const -16
        i32.and
        local.tee $4
        i32.const 88
        i32.add
        local.tee $5
        i32.store
        local.get $5
        i32.eqz
        br_if $block
        local.get $4
        i32.const 1656
        i32.add
        i32.const 0
        i32.store
        local.get $5
        i32.const 0
        i32.store
        local.get $4
        i32.const 184
        i32.add
        local.set $5
        i32.const 0
        local.set $6
        loop $loop
          block $block_2
            block $block_3
              local.get $6
              i32.const 23
              i32.eq
              br_if $block_3
              i32.const 0
              local.set $3
              local.get $4
              local.get $6
              i32.const 2
              i32.shl
              i32.add
              i32.const 92
              i32.add
              i32.const 0
              i32.store
              loop $loop_0
                local.get $3
                i32.const 64
                i32.eq
                br_if $block_2
                local.get $5
                local.get $3
                i32.add
                i32.const 0
                i32.store
                local.get $3
                i32.const 4
                i32.add
                local.set $3
                br $loop_0
              end ;; $loop_0
            end ;; $block_3
            i32.const 66879
            i32.const -4
            i32.and
            local.tee $3
            local.get $3
            i32.const 15
            i32.add
            i32.const -16
            i32.and
            i32.const 1660
            i32.add
            local.get $2
            call $84
            i32.const 0
            local.get $3
            i32.store offset=66856
            br $block_0
          end ;; $block_2
          local.get $5
          i32.const 64
          i32.add
          local.set $5
          local.get $6
          i32.const 1
          i32.add
          local.set $6
          br $loop
        end ;; $loop
      end ;; $block_0
      i32.const 344
      call $118
      local.set $3
      i32.const 512
      call $118
      i32.const 0
      i32.const 512
      call $145
      local.set $5
      local.get $3
      i32.eqz
      br_if $block
      local.get $3
      i32.const 0
      i32.store offset=32
      local.get $3
      i64.const -3221225472
      i64.store offset=24 align=4
      local.get $3
      i64.const 274877906954
      i64.store offset=16 align=4
      local.get $3
      i64.const 64
      i64.store offset=8 align=4
      local.get $3
      local.get $5
      i32.store
      local.get $3
      i32.const 48
      i32.add
      i32.const 2
      i32.store
      local.get $3
      i32.const 40
      i32.add
      i32.const 1
      i32.store
      local.get $3
      local.get $5
      i32.const 512
      i32.add
      i32.store offset=4
      local.get $0
      call $119
      i32.const 0
      local.get $3
      i32.store offset=66860
      local.get $3
      local.get $0
      i64.load
      local.tee $1
      i64.const 1073741823
      i64.and
      local.get $0
      i64.load offset=8
      local.get $1
      i64.const 30
      i64.shr_u
      i64.const 8589934591
      i64.and
      i64.const 59453308800
      i64.add
      local.get $1
      i64.const -1
      i64.gt_s
      select
      i64.const 1000000000
      i64.mul
      i64.add
      i64.const -6795364578871345152
      i64.add
      i64.store offset=56
      local.get $3
      call $96
      call $120
      call $121
      local.get $0
      i32.const 16
      i32.add
      global.set $24
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $118 (type $6)
    (param $0 i32)
    (result i32)
    i32.const 0
    i32.load offset=66856
    local.get $0
    call $88
    )
  
  (func $119 (type $0)
    (param $0 i32)
    (local $1 i64)
    (local $2 i64)
    (local $3 i64)
    (local $4 i64)
    call $116
    local.tee $1
    i64.const 1000000000
    i64.div_s
    local.tee $2
    i64.const -1000000000
    i64.mul
    local.get $1
    i64.add
    local.set $3
    block $block
      local.get $2
      i64.const 2682288000
      i64.add
      local.tee $4
      i64.const 8589934592
      i64.lt_u
      br_if $block
      local.get $0
      local.get $2
      i64.const 62135596800
      i64.add
      i64.store offset=8
      local.get $0
      local.get $3
      i64.const 32
      i64.shl
      i64.const 32
      i64.shr_s
      i64.store
      return
    end ;; $block
    local.get $0
    local.get $1
    i64.const 1
    i64.add
    i64.store offset=8
    local.get $0
    local.get $4
    i64.const 30
    i64.shl
    local.get $3
    i64.const 32
    i64.shl
    i64.const 32
    i64.shr_s
    i64.or
    i64.const -9223372036854775808
    i64.or
    i64.store
    )
  
  (func $120 (type $17)
    (param $0 i32)
    (param $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $2
    global.set $24
    i32.const 20
    call $97
    local.tee $3
    local.get $1
    i32.store offset=8
    local.get $3
    i32.const 3
    i32.store offset=4
    local.get $3
    i32.const 4
    i32.store
    local.get $3
    local.get $1
    local.get $3
    call $95
    i32.store offset=12
    i32.const 0
    memory.size
    i32.const 16
    i32.shl
    i32.store offset=66560
    local.get $2
    i32.const 0
    i32.store offset=12
    local.get $2
    i32.const 0
    i32.store offset=8
    local.get $2
    i32.const 12
    i32.add
    local.get $2
    i32.const 8
    i32.add
    call $30
    drop
    block $block
      block $block_0
        block $block_1
          local.get $2
          i32.load offset=12
          local.tee $4
          i32.eqz
          br_if $block_1
          local.get $4
          i32.const 1073741823
          i32.gt_u
          br_if $block_0
          local.get $4
          i32.const 2
          i32.shl
          call $97
          local.set $5
          local.get $2
          i32.load offset=8
          local.tee $6
          i32.const -1
          i32.le_s
          br_if $block_0
          local.get $6
          call $97
          local.set $7
          local.get $6
          i32.eqz
          br_if $block
          local.get $5
          local.get $7
          call $31
          drop
          local.get $4
          i32.const 536870912
          i32.ge_u
          br_if $block_0
          i32.const 0
          local.set $6
          local.get $4
          i32.const 3
          i32.shl
          call $97
          local.set $7
          i32.const 0
          local.get $4
          i32.store offset=66504
          i32.const 0
          local.get $4
          i32.store offset=66500
          i32.const 0
          local.get $7
          i32.store offset=66496
          loop $loop
            local.get $4
            local.get $6
            i32.eq
            br_if $block_1
            local.get $6
            local.get $4
            i32.ge_u
            br_if $block
            local.get $7
            i32.const 4
            i32.add
            local.get $5
            i32.load
            local.tee $8
            call $146
            i32.store
            local.get $7
            local.get $8
            i32.store
            local.get $5
            i32.const 4
            i32.add
            local.set $5
            local.get $7
            i32.const 8
            i32.add
            local.set $7
            local.get $6
            i32.const 1
            i32.add
            local.set $6
            br $loop
          end ;; $loop
        end ;; $block_1
        local.get $6
        local.get $1
        call $124
        local.get $3
        i32.const 0
        i32.store8 offset=16
        local.get $2
        i32.const 16
        i32.add
        global.set $24
        return
      end ;; $block_0
      call $99
      unreachable
    end ;; $block
    call $68
    unreachable
    )
  
  (func $121 (type $5)
    (local $0 i64)
    (local $1 i32)
    (local $2 i64)
    (local $3 i64)
    i64.const 0
    local.set $0
    loop $loop
      block $block
        block $block_0
          i32.const 0
          i32.load8_u offset=66840
          br_if $block_0
          block $block_1
            i32.const 0
            i32.load offset=66832
            i32.eqz
            br_if $block_1
            call $116
            local.set $0
            i32.const 0
            i32.load offset=66832
            local.tee $1
            i32.eqz
            br_if $block_1
            local.get $0
            i32.const 0
            i64.load offset=66848
            local.tee $2
            i64.sub
            local.get $1
            i64.load offset=8
            local.tee $3
            i64.lt_s
            br_if $block_1
            i32.const 0
            local.get $3
            local.get $2
            i64.add
            i64.store offset=66848
            i32.const 0
            local.get $1
            i32.load
            i32.store offset=66832
            local.get $1
            i32.const 0
            i32.store
            i32.const 66824
            local.get $1
            call $94
          end ;; $block_1
          call $93
          local.tee $1
          br_if $block
          i32.const 0
          i32.load offset=66832
          local.tee $1
          i32.eqz
          br_if $block_0
          local.get $1
          i64.load offset=8
          local.get $0
          i64.sub
          i32.const 0
          i64.load offset=66848
          i64.add
          call $32
        end ;; $block_0
        return
      end ;; $block
      local.get $1
      i32.load offset=16
      local.tee $1
      local.get $1
      i32.load
      call_indirect $22 (type $0)
      br $loop
    end ;; $loop
    )
  
  (func $122 (type $0)
    (param $0 i32)
    local.get $0
    call $104
    )
  
  (func $123 (type $0)
    (param $0 i32)
    (local $1 i32)
    i32.const 0
    i32.const 1
    i32.store8 offset=66840
    block $block
      local.get $0
      i32.load offset=8
      local.tee $1
      br_if $block
      call $87
      unreachable
    end ;; $block
    local.get $1
    local.get $0
    i32.load offset=12
    i32.store offset=16
    i32.const 66824
    local.get $1
    call $94
    local.get $0
    call $104
    )
  
  (func $124 (type $17)
    (param $0 i32)
    (param $1 i32)
    (local $2 i32)
    i32.const 44
    call $97
    local.tee $2
    local.get $1
    i32.store offset=32
    local.get $2
    i32.const 5
    i32.store offset=4
    local.get $2
    i32.const 6
    i32.store
    local.get $2
    local.get $1
    local.get $2
    call $95
    i32.store offset=36
    i32.const 66468
    i32.const 13
    call $101
    call $102
    local.get $2
    call $96
    call $140
    block $block
      local.get $1
      br_if $block
      call $87
      unreachable
    end ;; $block
    local.get $1
    local.get $2
    i32.const 41
    i32.add
    i32.store offset=4
    i32.const 66524
    local.get $2
    i32.const 42
    i32.add
    local.get $2
    i32.const 8
    i32.add
    local.get $2
    local.get $1
    call $113
    drop
    local.get $2
    i32.const 0
    i32.store8 offset=40
    )
  
  (func $125 (type $5)
    i32.const 0
    i32.const 0
    i32.store8 offset=66840
    call $121
    )
  
  (func $126 (type $6)
    (param $0 i32)
    (result i32)
    i32.const 0
    i32.load offset=66856
    local.get $0
    call $79
    )
  
  (func $127 (type $17)
    (param $0 i32)
    (param $1 i32)
    local.get $0
    i32.const 0
    i32.load offset=66856
    local.get $1
    call $79
    local.tee $1
    call $91
    i32.store offset=4
    local.get $0
    local.get $1
    i32.store
    )
  
  (func $128 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    i32.const 0
    i32.load offset=66856
    local.get $1
    local.get $0
    i32.mul
    call $88
    )
  
  (func $129 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    (local $2 i32)
    (local $3 i32)
    i32.const 0
    local.set $2
    i32.const 0
    i32.load offset=66856
    local.set $3
    local.get $0
    call $89
    local.set $0
    block $block
      block $block_0
        local.get $3
        local.get $1
        call $80
        local.tee $1
        i32.eqz
        br_if $block_0
        local.get $0
        i32.eqz
        br_if $block
        local.get $1
        i32.const 4
        i32.add
        local.tee $2
        local.get $0
        i32.const 4
        i32.add
        local.get $0
        i32.load
        i32.const -4
        i32.and
        call $144
        drop
        local.get $3
        local.get $0
        call $90
      end ;; $block_0
      local.get $2
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $130 (type $3)
    (param $0 i32)
    (param $1 i32)
    (result i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    block $block
      local.get $0
      i32.eqz
      br_if $block
      local.get $1
      call $92
      local.get $0
      i32.load offset=8
      i32.rem_u
      i32.const 3
      i32.shl
      local.get $0
      i32.load
      local.tee $2
      i32.add
      local.tee $3
      local.set $4
      block $block_0
        loop $loop
          local.get $4
          i32.load
          local.tee $5
          i32.eqz
          br_if $block_0
          block $block_1
            local.get $5
            local.get $1
            i32.ne
            br_if $block_1
            i32.const 1
            return
          end ;; $block_1
          local.get $4
          i32.const 8
          i32.add
          local.tee $4
          local.get $2
          local.get $4
          local.get $0
          i32.load offset=4
          i32.lt_u
          select
          local.tee $4
          local.get $3
          i32.ne
          br_if $loop
        end ;; $loop
      end ;; $block_0
      i32.const 0
      return
    end ;; $block
    call $87
    unreachable
    )
  
  (func $131 (type $16)
    (param $0 i32)
    (param $1 i32)
    (param $2 i32)
    (local $3 i32)
    (local $4 i32)
    block $block
      local.get $2
      i32.const 65
      i32.ge_s
      br_if $block
      local.get $1
      i32.load offset=4
      br_if $block
      local.get $1
      i32.const 1
      i32.store offset=4
      local.get $1
      i32.load offset=8
      local.tee $3
      i32.const 3
      i32.and
      br_if $block
      local.get $2
      i32.const 1
      i32.add
      local.set $4
      local.get $3
      local.get $1
      i32.const 12
      i32.add
      local.tee $2
      i32.add
      local.set $3
      loop $loop
        local.get $2
        local.get $3
        i32.ge_u
        br_if $block
        block $block_0
          local.get $0
          br_if $block_0
          call $87
          unreachable
        end ;; $block_0
        block $block_1
          local.get $2
          i32.load
          i32.const -12
          i32.add
          local.tee $1
          local.get $0
          i32.load offset=28
          i32.lt_u
          br_if $block_1
          local.get $1
          local.get $0
          i32.load offset=32
          i32.gt_u
          br_if $block_1
          local.get $0
          local.get $1
          call $130
          i32.const 1
          i32.and
          i32.eqz
          br_if $block_1
          local.get $0
          local.get $1
          local.get $4
          call $131
        end ;; $block_1
        local.get $2
        i32.const 4
        i32.add
        local.set $2
        br $loop
      end ;; $loop
    end ;; $block
    )
  
  (func $132 (type $17)
    (param $0 i32)
    (param $1 i32)
    block $block
      local.get $0
      br_if $block
      call $87
      unreachable
    end ;; $block
    block $block_0
      local.get $0
      i32.load offset=28
      local.get $1
      i32.gt_u
      br_if $block_0
      local.get $0
      i32.load offset=32
      local.get $1
      i32.lt_u
      br_if $block_0
      local.get $0
      local.get $1
      i32.const -12
      i32.add
      local.tee $1
      call $130
      i32.const 1
      i32.and
      i32.eqz
      br_if $block_0
      local.get $1
      i32.const 1
      i32.store offset=4
      return
    end ;; $block_0
    )
  
  (func $133 (type $16)
    (param $0 i32)
    (param $1 i32)
    (param $2 i32)
    block $block
      local.get $1
      i32.eqz
      br_if $block
      local.get $2
      local.get $1
      i32.le_u
      br_if $block
      local.get $2
      i32.const -4
      i32.and
      local.set $2
      local.get $1
      i32.const 3
      i32.add
      i32.const -4
      i32.and
      local.set $1
      loop $loop
        local.get $1
        local.get $2
        i32.ge_u
        br_if $block
        local.get $0
        local.get $1
        i32.load
        call $132
        local.get $1
        i32.const 4
        i32.add
        local.set $1
        br $loop
      end ;; $loop
    end ;; $block
    )
  
  (func $134 (type $5)
    (local $0 i32)
    (local $1 i32)
    i32.const 0
    local.set $0
    block $block
      loop $loop
        local.get $0
        i32.load
        local.tee $0
        i32.eqz
        br_if $block
        i32.const 0
        i32.load offset=66860
        local.get $0
        i32.const 8
        i32.add
        local.tee $1
        local.get $1
        local.get $0
        i32.load offset=4
        i32.const 2
        i32.shl
        i32.add
        call $133
        br $loop
      end ;; $loop
    end ;; $block
    )
  
  (func $135 (type $5)
    (local $0 i32)
    (local $1 i32)
    (local $2 i32)
    global.get $24
    i32.const 16
    i32.sub
    local.tee $0
    global.set $24
    i32.const 0
    i32.load offset=66860
    i32.const 65536
    i32.const 66876
    call $133
    i32.const 66832
    local.set $1
    block $block
      loop $loop
        local.get $1
        i32.load
        local.tee $2
        i32.eqz
        br_if $block
        local.get $2
        call $136
        local.get $1
        i32.load
        local.tee $1
        br_if $loop
      end ;; $loop
      call $87
      unreachable
    end ;; $block
    local.get $0
    i64.const 0
    i64.store offset=8 align=4
    block $block_0
      loop $loop_0
        i32.const 0
        i32.load offset=66824
        i32.eqz
        br_if $block_0
        call $93
        local.tee $1
        call $136
        local.get $0
        i32.const 8
        i32.add
        local.get $1
        call $94
        br $loop_0
      end ;; $loop_0
    end ;; $block_0
    i32.const 0
    local.get $0
    i64.load offset=8 align=4
    i64.store offset=66824
    local.get $0
    i32.const 16
    i32.add
    global.set $24
    )
  
  (func $136 (type $0)
    (param $0 i32)
    i32.const 0
    i32.load offset=66860
    local.get $0
    call $132
    )
  
  (func $137 (type $5)
    )
  
  (func $138 (type $0)
    (param $0 i32)
    local.get $0
    call $104
    )
  
  (func $139 (type $0)
    (param $0 i32)
    (local $1 i32)
    block $block
      local.get $0
      i32.load offset=32
      local.tee $1
      br_if $block
      call $87
      unreachable
    end ;; $block
    local.get $1
    local.get $0
    i32.load offset=36
    i32.store offset=16
    i32.const 66824
    local.get $1
    call $94
    local.get $0
    call $104
    )
  
  (func $140 (type $17)
    (param $0 i32)
    (param $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i64)
    (local $8 i64)
    (local $9 i64)
    (local $10 i64)
    (local $11 i64)
    (local $12 i32)
    (local $13 i32)
    (local $14 i64)
    (local $15 i64)
    (local $16 i64)
    (local $17 i64)
    global.get $24
    i32.const 96
    i32.sub
    local.tee $2
    global.set $24
    i32.const 16
    call $97
    local.tee $3
    local.get $1
    i32.store offset=8
    local.get $3
    i32.const 7
    i32.store offset=4
    local.get $3
    i32.const 8
    i32.store
    local.get $1
    local.get $3
    call $95
    drop
    local.get $2
    i32.const 88
    i32.add
    i32.const 152
    call $127
    block $block
      local.get $2
      i32.load offset=88
      local.tee $1
      i32.eqz
      br_if $block
      local.get $2
      i32.load offset=92
      local.set $4
      local.get $1
      i64.const 32
      i64.store offset=4 align=4
      local.get $1
      local.get $4
      i32.store offset=12
      local.get $1
      local.get $1
      i32.const 16
      i32.add
      local.tee $4
      i32.store
      local.get $4
      i32.eqz
      br_if $block
      local.get $4
      i64.const 0
      i64.store align=4
      i32.const 512
      call $126
      local.set $5
      block $block_0
        local.get $1
        i32.load
        local.tee $4
        i32.eqz
        br_if $block_0
        block $block_1
          block $block_2
            local.get $4
            i32.load
            local.tee $6
            local.get $1
            i32.load offset=4
            i32.ne
            br_if $block_2
            local.get $2
            i32.const 80
            i32.add
            local.get $6
            i32.const 2
            i32.shl
            i32.const 8
            i32.add
            call $127
            local.get $2
            i32.load offset=80
            local.set $6
            local.get $1
            local.get $2
            i32.load offset=84
            local.get $5
            call $91
            i32.add
            local.get $1
            i32.load offset=12
            i32.add
            i32.store offset=12
            local.get $6
            i32.eqz
            br_if $block
            local.get $6
            local.get $4
            i32.store offset=4
            local.get $6
            i32.const 1
            i32.store
            local.get $1
            local.get $6
            i32.store
            local.get $6
            local.get $5
            i32.store offset=8
            br $block_1
          end ;; $block_2
          local.get $1
          local.get $5
          call $91
          local.get $1
          i32.load offset=12
          i32.add
          i32.store offset=12
          local.get $4
          local.get $4
          i32.load
          i32.const 2
          i32.shl
          i32.add
          i32.const 8
          i32.add
          local.get $5
          i32.store
          local.get $4
          local.get $4
          i32.load
          i32.const 1
          i32.add
          i32.store
        end ;; $block_1
        local.get $1
        local.get $1
        i32.load offset=8
        i32.const 1
        i32.add
        i32.store offset=8
      end ;; $block_0
      block $block_3
        local.get $1
        i32.load
        local.tee $6
        i32.eqz
        br_if $block_3
        loop $loop
          local.get $6
          i32.eqz
          br_if $block_3
          local.get $6
          i32.load
          i32.const 2
          i32.shl
          local.get $6
          i32.const 8
          i32.add
          local.tee $1
          i32.add
          local.set $5
          block $block_4
            loop $loop_0
              local.get $1
              local.get $5
              i32.ge_u
              br_if $block_4
              local.get $1
              i32.load
              local.tee $4
              i32.eqz
              br_if $block_4
              local.get $4
              call $108
              local.get $1
              i32.const 4
              i32.add
              local.set $1
              br $loop_0
            end ;; $loop_0
          end ;; $block_4
          block $block_5
            local.get $6
            i32.load offset=4
            i32.eqz
            br_if $block_5
            local.get $6
            call $108
            local.get $6
            i32.load offset=4
            local.set $6
            br $loop
          end ;; $block_5
        end ;; $loop
        local.get $6
        i32.const -16
        i32.add
        call $108
      end ;; $block_3
      block $block_6
        i32.const 0
        i32.load offset=66872
        br_if $block_6
        i32.const 128
        call $97
        local.tee $1
        i32.const 10
        i32.store8
        i32.const 0
        local.get $1
        i32.store offset=66872
      end ;; $block_6
      local.get $2
      i32.const 64
      i32.add
      call $119
      local.get $2
      i64.load offset=64
      local.tee $7
      i64.const 1073741823
      i64.and
      local.get $2
      i64.load offset=72
      local.get $7
      i64.const 30
      i64.shr_u
      i64.const 8589934591
      i64.and
      i64.const 59453308800
      i64.add
      local.get $7
      i64.const -1
      i64.gt_s
      select
      i64.const 1000000000
      i64.mul
      i64.add
      i64.const -6795364578871345152
      i64.add
      call $112
      call $102
      i32.const 0
      i32.load offset=66860
      local.set $4
      local.get $2
      i32.const 48
      i32.add
      call $119
      local.get $4
      i32.eqz
      br_if $block
      local.get $2
      i64.load offset=48
      local.tee $7
      i64.const 30
      i64.shr_u
      i64.const 8589934591
      i64.and
      local.set $8
      local.get $2
      i64.load offset=56
      local.set $9
      local.get $4
      i32.const 64
      i32.add
      local.tee $1
      local.get $1
      i64.load
      i64.const 1
      i64.add
      i64.store
      block $block_7
        block $block_8
          block $block_9
            local.get $4
            i32.const 48
            i32.add
            i32.load
            br_table
              $block_7 $block_9 $block_8
              $block_9 ;; default
          end ;; $block_9
          call $135
          br $block_7
        end ;; $block_8
        call $134
      end ;; $block_7
      local.get $7
      i64.const -1
      i64.gt_s
      local.set $1
      local.get $8
      i64.const 59453308800
      i64.add
      local.set $8
      block $block_10
        block $block_11
          block $block_12
            local.get $4
            i32.const 40
            i32.add
            i32.load
            br_table
              $block_10 $block_11 $block_12
              $block_11 ;; default
          end ;; $block_12
          call $134
          br $block_10
        end ;; $block_11
        call $135
      end ;; $block_10
      local.get $9
      local.get $8
      local.get $1
      select
      local.set $10
      local.get $7
      i64.const 1073741823
      i64.and
      local.set $11
      local.get $2
      i32.const 32
      i32.add
      call $119
      local.get $2
      i64.load offset=40
      local.set $8
      local.get $2
      i64.load offset=32
      local.set $7
      local.get $4
      i32.const 336
      i32.add
      i64.const 0
      i64.store
      local.get $4
      i32.const 328
      i32.add
      i64.const 0
      i64.store
      local.get $4
      i32.load offset=8
      i32.const 3
      i32.shl
      local.get $4
      i32.load
      local.tee $12
      i32.add
      local.set $13
      local.get $7
      i64.const 30
      i64.shr_u
      i64.const 8589934591
      i64.and
      i64.const 59453308800
      i64.add
      local.set $9
      loop $loop_1
        block $block_13
          block $block_14
            local.get $12
            local.get $13
            i32.ge_u
            br_if $block_14
            local.get $12
            i32.load
            local.tee $1
            i32.eqz
            br_if $block_13
            local.get $1
            i32.load offset=8
            local.tee $5
            i32.const 3
            i32.and
            br_if $block_13
            local.get $5
            local.get $1
            i32.const 12
            i32.add
            local.tee $1
            i32.add
            local.set $6
            loop $loop_2
              local.get $1
              local.get $6
              i32.ge_u
              br_if $block_13
              block $block_15
                local.get $1
                i32.load
                i32.const -12
                i32.add
                local.tee $5
                local.get $4
                i32.load offset=28
                i32.lt_u
                br_if $block_15
                local.get $5
                local.get $4
                i32.load offset=32
                i32.gt_u
                br_if $block_15
                local.get $4
                local.get $5
                call $130
                i32.const 1
                i32.and
                i32.eqz
                br_if $block_15
                local.get $4
                local.get $5
                i32.const 0
                call $131
              end ;; $block_15
              local.get $1
              i32.const 4
              i32.add
              local.set $1
              br $loop_2
            end ;; $loop_2
          end ;; $block_14
          local.get $8
          local.get $9
          local.get $7
          i64.const -1
          i64.gt_s
          select
          local.set $14
          local.get $7
          i64.const 1073741823
          i64.and
          local.set $15
          local.get $2
          i32.const 16
          i32.add
          call $119
          local.get $4
          i32.load offset=8
          i32.const 3
          i32.shl
          local.get $4
          i32.load
          local.tee $5
          i32.add
          local.set $13
          i32.const -1
          local.set $6
          local.get $2
          i64.load offset=24
          local.set $16
          local.get $2
          i64.load offset=16
          local.set $9
          i32.const 0
          local.set $12
          loop $loop_3
            block $block_16
              block $block_17
                local.get $5
                local.get $13
                i32.ge_u
                br_if $block_17
                local.get $5
                i32.load
                local.tee $1
                i32.eqz
                br_if $block_16
                block $block_18
                  local.get $1
                  i32.load offset=4
                  i32.eqz
                  br_if $block_18
                  local.get $1
                  i32.const 0
                  i32.store offset=4
                  local.get $1
                  local.get $12
                  local.get $1
                  local.get $12
                  i32.gt_u
                  select
                  local.set $12
                  local.get $1
                  local.get $6
                  local.get $1
                  local.get $6
                  i32.lt_u
                  select
                  local.set $6
                  br $block_16
                end ;; $block_18
                local.get $4
                local.get $4
                i32.load offset=344
                local.get $1
                i32.load
                i32.sub
                i32.store offset=344
                local.get $1
                i64.load32_u
                local.set $7
                local.get $4
                local.get $4
                i64.load offset=72
                i64.const -1
                i64.add
                i64.store offset=72
                local.get $4
                local.get $4
                i64.load offset=328
                i64.const 1
                i64.add
                i64.store offset=328
                local.get $4
                local.get $7
                local.get $4
                i64.load offset=336
                i64.add
                i64.store offset=336
                local.get $1
                i64.load32_u offset=8
                local.set $7
                local.get $1
                i64.load32_u
                local.set $8
                i32.const 66192
                i32.const 8
                call $101
                i32.const 32
                call $103
                local.get $1
                i32.const 12
                i32.add
                i64.extend_i32_u
                call $105
                i32.const 32
                call $103
                i32.const 66200
                i32.const 4
                call $101
                i32.const 32
                call $103
                local.get $8
                call $105
                i32.const 32
                call $103
                i32.const 66204
                i32.const 6
                call $101
                i32.const 32
                call $103
                local.get $7
                call $105
                call $102
                local.get $1
                call $108
                local.get $4
                local.get $1
                call $106
                drop
                br $block_16
              end ;; $block_17
              local.get $4
              local.get $12
              i32.store offset=32
              local.get $4
              local.get $6
              i32.store offset=28
              local.get $2
              call $119
              local.get $2
              i64.load offset=8
              local.set $17
              local.get $2
              i64.load
              local.set $7
              local.get $4
              i32.const 296
              i32.add
              local.get $15
              local.get $14
              i64.const 1000000000
              i64.mul
              i64.add
              i64.const -6795364578871345152
              i64.add
              local.tee $14
              local.get $10
              i64.const -1000000000
              i64.mul
              local.get $11
              i64.sub
              i64.const 6795364578871345152
              i64.add
              local.tee $10
              i64.add
              i64.store
              local.get $4
              i32.const 120
              i32.add
              local.tee $1
              local.get $1
              i64.load
              local.get $4
              i64.load offset=336
              i64.add
              i64.store
              local.get $4
              i32.const 112
              i32.add
              local.tee $1
              local.get $1
              i64.load
              local.get $4
              i64.load offset=328
              i64.add
              i64.store
              local.get $4
              i32.const 304
              i32.add
              local.get $9
              i64.const 1073741823
              i64.and
              local.get $16
              local.get $9
              i64.const 30
              i64.shr_u
              i64.const 8589934591
              i64.and
              i64.const 59453308800
              i64.add
              local.get $9
              i64.const -1
              i64.gt_s
              select
              i64.const 1000000000
              i64.mul
              i64.add
              i64.const -6795364578871345152
              i64.add
              local.tee $8
              local.get $14
              i64.sub
              i64.store
              local.get $4
              i32.const 312
              i32.add
              local.get $7
              i64.const 1073741823
              i64.and
              local.get $8
              i64.sub
              local.get $17
              local.get $7
              i64.const 30
              i64.shr_u
              i64.const 8589934591
              i64.and
              i64.const 59453308800
              i64.add
              local.get $7
              i64.const -1
              i64.gt_s
              select
              i64.const 1000000000
              i64.mul
              i64.add
              i64.const -6795364578871345152
              i64.add
              local.tee $7
              i64.store
              local.get $4
              i32.const 128
              i32.add
              local.tee $1
              local.get $7
              local.get $1
              i64.load
              i64.add
              i64.store
              local.get $4
              i32.const 320
              i32.add
              local.get $8
              local.get $10
              i64.add
              local.get $7
              i64.add
              local.tee $7
              i64.store
              local.get $4
              i32.const 264
              i32.add
              local.tee $1
              local.get $1
              i64.load
              local.get $7
              i64.add
              i64.store
              i32.const 0
              i32.load offset=66860
              local.tee $1
              i32.eqz
              br_if $block
              local.get $3
              i32.load offset=8
              local.set $4
              i32.const 66244
              i32.const 8
              call $101
              call $102
              local.get $1
              i32.const 72
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66252
              i32.const 10
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 344
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66262
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 96
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66277
              i32.const 11
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 80
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66288
              i32.const 12
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 104
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66300
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 120
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66315
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 88
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66330
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 328
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66345
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 336
              i32.add
              i64.load32_u
              local.set $7
              i32.const 66360
              i32.const 19
              call $101
              i32.const 32
              call $103
              local.get $7
              call $105
              call $102
              local.get $1
              i32.const 296
              i32.add
              i64.load
              local.set $7
              i32.const 66379
              i32.const 18
              call $101
              i32.const 32
              call $103
              local.get $7
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $1
              i32.const 304
              i32.add
              i64.load
              local.set $7
              i32.const 66397
              i32.const 18
              call $101
              i32.const 32
              call $103
              local.get $7
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $1
              i32.const 312
              i32.add
              i64.load
              local.set $7
              i32.const 66415
              i32.const 18
              call $101
              i32.const 32
              call $103
              local.get $7
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $1
              i32.const 320
              i32.add
              i64.load
              local.set $7
              i32.const 66433
              i32.const 16
              call $101
              i32.const 32
              call $103
              local.get $7
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $4
              i32.eqz
              br_if $block
              local.get $3
              i32.load offset=8
              i64.const 1000000000
              i64.store offset=8
              call $116
              local.set $7
              block $block_19
                i32.const 0
                i32.load offset=66832
                local.tee $4
                br_if $block_19
                i32.const 0
                local.get $7
                i64.store offset=66848
              end ;; $block_19
              i32.const 66832
              local.set $1
              block $block_20
                loop $loop_4
                  local.get $4
                  i32.eqz
                  br_if $block_20
                  block $block_21
                    local.get $3
                    i32.load offset=8
                    local.tee $5
                    i64.load offset=8
                    local.tee $7
                    local.get $4
                    i64.load offset=8
                    local.tee $8
                    i64.lt_u
                    br_if $block_21
                    local.get $5
                    local.get $7
                    local.get $8
                    i64.sub
                    i64.store offset=8
                    local.get $1
                    i32.load
                    local.tee $1
                    i32.eqz
                    br_if $block
                    local.get $1
                    i32.load
                    local.set $4
                    br $loop_4
                  end ;; $block_21
                end ;; $loop_4
                local.get $4
                local.get $8
                local.get $7
                i64.sub
                i64.store offset=8
              end ;; $block_20
              local.get $3
              i32.load offset=8
              local.tee $4
              local.get $1
              i32.load
              i32.store
              local.get $1
              local.get $4
              i32.store
              local.get $3
              i32.const 0
              i32.store8 offset=12
              local.get $2
              i32.const 96
              i32.add
              global.set $24
              return
            end ;; $block_16
            local.get $5
            i32.const 8
            i32.add
            local.set $5
            br $loop_3
          end ;; $loop_3
        end ;; $block_13
        local.get $12
        i32.const 8
        i32.add
        local.set $12
        br $loop_1
      end ;; $loop_1
    end ;; $block
    call $87
    unreachable
    )
  
  (func $141 (type $0)
    (param $0 i32)
    local.get $0
    call $104
    )
  
  (func $142 (type $0)
    (param $0 i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i64)
    (local $7 i64)
    (local $8 i64)
    (local $9 i64)
    (local $10 i64)
    (local $11 i32)
    (local $12 i32)
    (local $13 i64)
    (local $14 i64)
    (local $15 i64)
    (local $16 i64)
    global.get $24
    i32.const 96
    i32.sub
    local.tee $1
    global.set $24
    local.get $1
    i32.const 88
    i32.add
    i32.const 152
    call $127
    block $block
      local.get $1
      i32.load offset=88
      local.tee $2
      i32.eqz
      br_if $block
      local.get $1
      i32.load offset=92
      local.set $3
      local.get $2
      i64.const 32
      i64.store offset=4 align=4
      local.get $2
      local.get $3
      i32.store offset=12
      local.get $2
      local.get $2
      i32.const 16
      i32.add
      local.tee $3
      i32.store
      local.get $3
      i32.eqz
      br_if $block
      local.get $3
      i64.const 0
      i64.store align=4
      i32.const 512
      call $126
      local.set $4
      block $block_0
        local.get $2
        i32.load
        local.tee $3
        i32.eqz
        br_if $block_0
        block $block_1
          block $block_2
            local.get $3
            i32.load
            local.tee $5
            local.get $2
            i32.load offset=4
            i32.ne
            br_if $block_2
            local.get $1
            i32.const 80
            i32.add
            local.get $5
            i32.const 2
            i32.shl
            i32.const 8
            i32.add
            call $127
            local.get $1
            i32.load offset=80
            local.set $5
            local.get $2
            local.get $1
            i32.load offset=84
            local.get $4
            call $91
            i32.add
            local.get $2
            i32.load offset=12
            i32.add
            i32.store offset=12
            local.get $5
            i32.eqz
            br_if $block
            local.get $5
            local.get $3
            i32.store offset=4
            local.get $5
            i32.const 1
            i32.store
            local.get $2
            local.get $5
            i32.store
            local.get $5
            local.get $4
            i32.store offset=8
            br $block_1
          end ;; $block_2
          local.get $2
          local.get $4
          call $91
          local.get $2
          i32.load offset=12
          i32.add
          i32.store offset=12
          local.get $3
          local.get $3
          i32.load
          i32.const 2
          i32.shl
          i32.add
          i32.const 8
          i32.add
          local.get $4
          i32.store
          local.get $3
          local.get $3
          i32.load
          i32.const 1
          i32.add
          i32.store
        end ;; $block_1
        local.get $2
        local.get $2
        i32.load offset=8
        i32.const 1
        i32.add
        i32.store offset=8
      end ;; $block_0
      block $block_3
        local.get $2
        i32.load
        local.tee $5
        i32.eqz
        br_if $block_3
        loop $loop
          local.get $5
          i32.eqz
          br_if $block_3
          local.get $5
          i32.load
          i32.const 2
          i32.shl
          local.get $5
          i32.const 8
          i32.add
          local.tee $2
          i32.add
          local.set $4
          block $block_4
            loop $loop_0
              local.get $2
              local.get $4
              i32.ge_u
              br_if $block_4
              local.get $2
              i32.load
              local.tee $3
              i32.eqz
              br_if $block_4
              local.get $3
              call $108
              local.get $2
              i32.const 4
              i32.add
              local.set $2
              br $loop_0
            end ;; $loop_0
          end ;; $block_4
          block $block_5
            local.get $5
            i32.load offset=4
            i32.eqz
            br_if $block_5
            local.get $5
            call $108
            local.get $5
            i32.load offset=4
            local.set $5
            br $loop
          end ;; $block_5
        end ;; $loop
        local.get $5
        i32.const -16
        i32.add
        call $108
      end ;; $block_3
      block $block_6
        i32.const 0
        i32.load offset=66872
        br_if $block_6
        i32.const 128
        call $97
        local.tee $2
        i32.const 10
        i32.store8
        i32.const 0
        local.get $2
        i32.store offset=66872
      end ;; $block_6
      local.get $1
      i32.const 64
      i32.add
      call $119
      local.get $1
      i64.load offset=64
      local.tee $6
      i64.const 1073741823
      i64.and
      local.get $1
      i64.load offset=72
      local.get $6
      i64.const 30
      i64.shr_u
      i64.const 8589934591
      i64.and
      i64.const 59453308800
      i64.add
      local.get $6
      i64.const -1
      i64.gt_s
      select
      i64.const 1000000000
      i64.mul
      i64.add
      i64.const -6795364578871345152
      i64.add
      call $112
      call $102
      i32.const 0
      i32.load offset=66860
      local.set $3
      local.get $1
      i32.const 48
      i32.add
      call $119
      local.get $3
      i32.eqz
      br_if $block
      local.get $1
      i64.load offset=48
      local.tee $6
      i64.const 30
      i64.shr_u
      i64.const 8589934591
      i64.and
      local.set $7
      local.get $1
      i64.load offset=56
      local.set $8
      local.get $3
      i32.const 64
      i32.add
      local.tee $2
      local.get $2
      i64.load
      i64.const 1
      i64.add
      i64.store
      block $block_7
        block $block_8
          block $block_9
            local.get $3
            i32.const 48
            i32.add
            i32.load
            br_table
              $block_7 $block_9 $block_8
              $block_9 ;; default
          end ;; $block_9
          call $135
          br $block_7
        end ;; $block_8
        call $134
      end ;; $block_7
      local.get $6
      i64.const -1
      i64.gt_s
      local.set $2
      local.get $7
      i64.const 59453308800
      i64.add
      local.set $7
      block $block_10
        block $block_11
          block $block_12
            local.get $3
            i32.const 40
            i32.add
            i32.load
            br_table
              $block_10 $block_11 $block_12
              $block_11 ;; default
          end ;; $block_12
          call $134
          br $block_10
        end ;; $block_11
        call $135
      end ;; $block_10
      local.get $8
      local.get $7
      local.get $2
      select
      local.set $9
      local.get $6
      i64.const 1073741823
      i64.and
      local.set $10
      local.get $1
      i32.const 32
      i32.add
      call $119
      local.get $1
      i64.load offset=40
      local.set $7
      local.get $1
      i64.load offset=32
      local.set $6
      local.get $3
      i32.const 336
      i32.add
      i64.const 0
      i64.store
      local.get $3
      i32.const 328
      i32.add
      i64.const 0
      i64.store
      local.get $3
      i32.load offset=8
      i32.const 3
      i32.shl
      local.get $3
      i32.load
      local.tee $11
      i32.add
      local.set $12
      local.get $6
      i64.const 30
      i64.shr_u
      i64.const 8589934591
      i64.and
      i64.const 59453308800
      i64.add
      local.set $8
      loop $loop_1
        block $block_13
          block $block_14
            local.get $11
            local.get $12
            i32.ge_u
            br_if $block_14
            local.get $11
            i32.load
            local.tee $2
            i32.eqz
            br_if $block_13
            local.get $2
            i32.load offset=8
            local.tee $4
            i32.const 3
            i32.and
            br_if $block_13
            local.get $4
            local.get $2
            i32.const 12
            i32.add
            local.tee $2
            i32.add
            local.set $5
            loop $loop_2
              local.get $2
              local.get $5
              i32.ge_u
              br_if $block_13
              block $block_15
                local.get $2
                i32.load
                i32.const -12
                i32.add
                local.tee $4
                local.get $3
                i32.load offset=28
                i32.lt_u
                br_if $block_15
                local.get $4
                local.get $3
                i32.load offset=32
                i32.gt_u
                br_if $block_15
                local.get $3
                local.get $4
                call $130
                i32.const 1
                i32.and
                i32.eqz
                br_if $block_15
                local.get $3
                local.get $4
                i32.const 0
                call $131
              end ;; $block_15
              local.get $2
              i32.const 4
              i32.add
              local.set $2
              br $loop_2
            end ;; $loop_2
          end ;; $block_14
          local.get $7
          local.get $8
          local.get $6
          i64.const -1
          i64.gt_s
          select
          local.set $13
          local.get $6
          i64.const 1073741823
          i64.and
          local.set $14
          local.get $1
          i32.const 16
          i32.add
          call $119
          local.get $3
          i32.load offset=8
          i32.const 3
          i32.shl
          local.get $3
          i32.load
          local.tee $4
          i32.add
          local.set $12
          i32.const -1
          local.set $5
          local.get $1
          i64.load offset=24
          local.set $15
          local.get $1
          i64.load offset=16
          local.set $8
          i32.const 0
          local.set $11
          loop $loop_3
            block $block_16
              block $block_17
                local.get $4
                local.get $12
                i32.ge_u
                br_if $block_17
                local.get $4
                i32.load
                local.tee $2
                i32.eqz
                br_if $block_16
                block $block_18
                  local.get $2
                  i32.load offset=4
                  i32.eqz
                  br_if $block_18
                  local.get $2
                  i32.const 0
                  i32.store offset=4
                  local.get $2
                  local.get $11
                  local.get $2
                  local.get $11
                  i32.gt_u
                  select
                  local.set $11
                  local.get $2
                  local.get $5
                  local.get $2
                  local.get $5
                  i32.lt_u
                  select
                  local.set $5
                  br $block_16
                end ;; $block_18
                local.get $3
                local.get $3
                i32.load offset=344
                local.get $2
                i32.load
                i32.sub
                i32.store offset=344
                local.get $2
                i64.load32_u
                local.set $6
                local.get $3
                local.get $3
                i64.load offset=72
                i64.const -1
                i64.add
                i64.store offset=72
                local.get $3
                local.get $3
                i64.load offset=328
                i64.const 1
                i64.add
                i64.store offset=328
                local.get $3
                local.get $6
                local.get $3
                i64.load offset=336
                i64.add
                i64.store offset=336
                local.get $2
                i64.load32_u offset=8
                local.set $6
                local.get $2
                i64.load32_u
                local.set $7
                i32.const 66192
                i32.const 8
                call $101
                i32.const 32
                call $103
                local.get $2
                i32.const 12
                i32.add
                i64.extend_i32_u
                call $105
                i32.const 32
                call $103
                i32.const 66200
                i32.const 4
                call $101
                i32.const 32
                call $103
                local.get $7
                call $105
                i32.const 32
                call $103
                i32.const 66204
                i32.const 6
                call $101
                i32.const 32
                call $103
                local.get $6
                call $105
                call $102
                local.get $2
                call $108
                local.get $3
                local.get $2
                call $106
                drop
                br $block_16
              end ;; $block_17
              local.get $3
              local.get $11
              i32.store offset=32
              local.get $3
              local.get $5
              i32.store offset=28
              local.get $1
              call $119
              local.get $1
              i64.load offset=8
              local.set $16
              local.get $1
              i64.load
              local.set $6
              local.get $3
              i32.const 296
              i32.add
              local.get $14
              local.get $13
              i64.const 1000000000
              i64.mul
              i64.add
              i64.const -6795364578871345152
              i64.add
              local.tee $13
              local.get $9
              i64.const -1000000000
              i64.mul
              local.get $10
              i64.sub
              i64.const 6795364578871345152
              i64.add
              local.tee $9
              i64.add
              i64.store
              local.get $3
              i32.const 120
              i32.add
              local.tee $2
              local.get $2
              i64.load
              local.get $3
              i64.load offset=336
              i64.add
              i64.store
              local.get $3
              i32.const 112
              i32.add
              local.tee $2
              local.get $2
              i64.load
              local.get $3
              i64.load offset=328
              i64.add
              i64.store
              local.get $3
              i32.const 304
              i32.add
              local.get $8
              i64.const 1073741823
              i64.and
              local.get $15
              local.get $8
              i64.const 30
              i64.shr_u
              i64.const 8589934591
              i64.and
              i64.const 59453308800
              i64.add
              local.get $8
              i64.const -1
              i64.gt_s
              select
              i64.const 1000000000
              i64.mul
              i64.add
              i64.const -6795364578871345152
              i64.add
              local.tee $7
              local.get $13
              i64.sub
              i64.store
              local.get $3
              i32.const 312
              i32.add
              local.get $6
              i64.const 1073741823
              i64.and
              local.get $7
              i64.sub
              local.get $16
              local.get $6
              i64.const 30
              i64.shr_u
              i64.const 8589934591
              i64.and
              i64.const 59453308800
              i64.add
              local.get $6
              i64.const -1
              i64.gt_s
              select
              i64.const 1000000000
              i64.mul
              i64.add
              i64.const -6795364578871345152
              i64.add
              local.tee $6
              i64.store
              local.get $3
              i32.const 128
              i32.add
              local.tee $2
              local.get $6
              local.get $2
              i64.load
              i64.add
              i64.store
              local.get $3
              i32.const 320
              i32.add
              local.get $7
              local.get $9
              i64.add
              local.get $6
              i64.add
              local.tee $6
              i64.store
              local.get $3
              i32.const 264
              i32.add
              local.tee $2
              local.get $2
              i64.load
              local.get $6
              i64.add
              i64.store
              i32.const 0
              i32.load offset=66860
              local.tee $2
              i32.eqz
              br_if $block
              local.get $0
              i32.load offset=8
              local.set $3
              i32.const 66244
              i32.const 8
              call $101
              call $102
              local.get $2
              i32.const 72
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66252
              i32.const 10
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 344
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66262
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 96
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66277
              i32.const 11
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 80
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66288
              i32.const 12
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 104
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66300
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 120
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66315
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 88
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66330
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 328
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66345
              i32.const 15
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 336
              i32.add
              i64.load32_u
              local.set $6
              i32.const 66360
              i32.const 19
              call $101
              i32.const 32
              call $103
              local.get $6
              call $105
              call $102
              local.get $2
              i32.const 296
              i32.add
              i64.load
              local.set $6
              i32.const 66379
              i32.const 18
              call $101
              i32.const 32
              call $103
              local.get $6
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $2
              i32.const 304
              i32.add
              i64.load
              local.set $6
              i32.const 66397
              i32.const 18
              call $101
              i32.const 32
              call $103
              local.get $6
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $2
              i32.const 312
              i32.add
              i64.load
              local.set $6
              i32.const 66415
              i32.const 18
              call $101
              i32.const 32
              call $103
              local.get $6
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $2
              i32.const 320
              i32.add
              i64.load
              local.set $6
              i32.const 66433
              i32.const 16
              call $101
              i32.const 32
              call $103
              local.get $6
              i64.const 1000
              i64.div_s
              call $112
              i32.const 32
              call $103
              i32.const 66449
              i32.const 3
              call $101
              call $102
              local.get $3
              i32.eqz
              br_if $block
              local.get $3
              i64.const 1000000000
              i64.store offset=8
              call $116
              local.set $6
              block $block_19
                i32.const 0
                i32.load offset=66832
                local.tee $3
                br_if $block_19
                i32.const 0
                local.get $6
                i64.store offset=66848
              end ;; $block_19
              local.get $0
              i32.load offset=8
              local.set $4
              i32.const 66832
              local.set $2
              block $block_20
                loop $loop_4
                  local.get $3
                  i32.eqz
                  br_if $block_20
                  block $block_21
                    local.get $4
                    i64.load offset=8
                    local.tee $6
                    local.get $3
                    i64.load offset=8
                    local.tee $7
                    i64.lt_u
                    br_if $block_21
                    local.get $4
                    local.get $6
                    local.get $7
                    i64.sub
                    i64.store offset=8
                    local.get $2
                    i32.load
                    local.tee $2
                    i32.eqz
                    br_if $block
                    local.get $2
                    i32.load
                    local.set $3
                    br $loop_4
                  end ;; $block_21
                end ;; $loop_4
                local.get $3
                local.get $7
                local.get $6
                i64.sub
                i64.store offset=8
              end ;; $block_20
              local.get $4
              local.get $2
              i32.load
              i32.store
              local.get $0
              i32.const 0
              i32.store8 offset=12
              local.get $2
              local.get $4
              i32.store
              local.get $1
              i32.const 96
              i32.add
              global.set $24
              return
            end ;; $block_16
            local.get $4
            i32.const 8
            i32.add
            local.set $4
            br $loop_3
          end ;; $loop_3
        end ;; $block_13
        local.get $11
        i32.const 8
        i32.add
        local.set $11
        br $loop_1
      end ;; $loop_1
    end ;; $block
    call $87
    unreachable
    )
  
  (func $143 (type $0)
    (param $0 i32)
    )
  
  (func $144 (type $19)
    (param $0 i32)
    (param $1 i32)
    (param $2 i32)
    (result i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i32)
    (local $7 i32)
    (local $8 i32)
    block $block
      block $block_0
        local.get $2
        i32.eqz
        br_if $block_0
        local.get $1
        i32.const 3
        i32.and
        i32.eqz
        br_if $block_0
        local.get $0
        local.set $3
        loop $loop
          local.get $3
          local.get $1
          i32.load8_u
          i32.store8
          local.get $2
          i32.const -1
          i32.add
          local.set $4
          local.get $3
          i32.const 1
          i32.add
          local.set $3
          local.get $1
          i32.const 1
          i32.add
          local.set $1
          local.get $2
          i32.const 1
          i32.eq
          br_if $block
          local.get $4
          local.set $2
          local.get $1
          i32.const 3
          i32.and
          br_if $loop
          br $block
        end ;; $loop
      end ;; $block_0
      local.get $2
      local.set $4
      local.get $0
      local.set $3
    end ;; $block
    block $block_1
      block $block_2
        local.get $3
        i32.const 3
        i32.and
        local.tee $2
        br_if $block_2
        block $block_3
          local.get $4
          i32.const 16
          i32.lt_u
          br_if $block_3
          loop $loop_0
            local.get $3
            local.get $1
            i32.load
            i32.store
            local.get $3
            i32.const 4
            i32.add
            local.get $1
            i32.const 4
            i32.add
            i32.load
            i32.store
            local.get $3
            i32.const 8
            i32.add
            local.get $1
            i32.const 8
            i32.add
            i32.load
            i32.store
            local.get $3
            i32.const 12
            i32.add
            local.get $1
            i32.const 12
            i32.add
            i32.load
            i32.store
            local.get $3
            i32.const 16
            i32.add
            local.set $3
            local.get $1
            i32.const 16
            i32.add
            local.set $1
            local.get $4
            i32.const -16
            i32.add
            local.tee $4
            i32.const 15
            i32.gt_u
            br_if $loop_0
          end ;; $loop_0
        end ;; $block_3
        block $block_4
          local.get $4
          i32.const 8
          i32.and
          i32.eqz
          br_if $block_4
          local.get $3
          local.get $1
          i64.load align=4
          i64.store align=4
          local.get $1
          i32.const 8
          i32.add
          local.set $1
          local.get $3
          i32.const 8
          i32.add
          local.set $3
        end ;; $block_4
        block $block_5
          local.get $4
          i32.const 4
          i32.and
          i32.eqz
          br_if $block_5
          local.get $3
          local.get $1
          i32.load
          i32.store
          local.get $1
          i32.const 4
          i32.add
          local.set $1
          local.get $3
          i32.const 4
          i32.add
          local.set $3
        end ;; $block_5
        block $block_6
          local.get $4
          i32.const 2
          i32.and
          i32.eqz
          br_if $block_6
          local.get $3
          local.get $1
          i32.load8_u
          i32.store8
          local.get $3
          local.get $1
          i32.load8_u offset=1
          i32.store8 offset=1
          local.get $3
          i32.const 2
          i32.add
          local.set $3
          local.get $1
          i32.const 2
          i32.add
          local.set $1
        end ;; $block_6
        local.get $4
        i32.const 1
        i32.and
        i32.eqz
        br_if $block_1
        local.get $3
        local.get $1
        i32.load8_u
        i32.store8
        local.get $0
        return
      end ;; $block_2
      block $block_7
        local.get $4
        i32.const 32
        i32.lt_u
        br_if $block_7
        block $block_8
          block $block_9
            block $block_10
              local.get $2
              i32.const -1
              i32.add
              br_table
                $block_10 $block_9 $block_8
                $block_7 ;; default
            end ;; $block_10
            local.get $3
            local.get $1
            i32.load8_u offset=1
            i32.store8 offset=1
            local.get $3
            local.get $1
            i32.load
            local.tee $5
            i32.store8
            local.get $3
            local.get $1
            i32.load8_u offset=2
            i32.store8 offset=2
            local.get $4
            i32.const -3
            i32.add
            local.set $4
            local.get $3
            i32.const 3
            i32.add
            local.set $6
            i32.const 0
            local.set $2
            loop $loop_1
              local.get $6
              local.get $2
              i32.add
              local.tee $3
              local.get $1
              local.get $2
              i32.add
              local.tee $7
              i32.const 4
              i32.add
              i32.load
              local.tee $8
              i32.const 8
              i32.shl
              local.get $5
              i32.const 24
              i32.shr_u
              i32.or
              i32.store
              local.get $3
              i32.const 4
              i32.add
              local.get $7
              i32.const 8
              i32.add
              i32.load
              local.tee $5
              i32.const 8
              i32.shl
              local.get $8
              i32.const 24
              i32.shr_u
              i32.or
              i32.store
              local.get $3
              i32.const 8
              i32.add
              local.get $7
              i32.const 12
              i32.add
              i32.load
              local.tee $8
              i32.const 8
              i32.shl
              local.get $5
              i32.const 24
              i32.shr_u
              i32.or
              i32.store
              local.get $3
              i32.const 12
              i32.add
              local.get $7
              i32.const 16
              i32.add
              i32.load
              local.tee $5
              i32.const 8
              i32.shl
              local.get $8
              i32.const 24
              i32.shr_u
              i32.or
              i32.store
              local.get $2
              i32.const 16
              i32.add
              local.set $2
              local.get $4
              i32.const -16
              i32.add
              local.tee $4
              i32.const 16
              i32.gt_u
              br_if $loop_1
            end ;; $loop_1
            local.get $6
            local.get $2
            i32.add
            local.set $3
            local.get $1
            local.get $2
            i32.add
            i32.const 3
            i32.add
            local.set $1
            br $block_7
          end ;; $block_9
          local.get $3
          local.get $1
          i32.load
          local.tee $5
          i32.store8
          local.get $3
          local.get $1
          i32.load8_u offset=1
          i32.store8 offset=1
          local.get $4
          i32.const -2
          i32.add
          local.set $4
          local.get $3
          i32.const 2
          i32.add
          local.set $6
          i32.const 0
          local.set $2
          loop $loop_2
            local.get $6
            local.get $2
            i32.add
            local.tee $3
            local.get $1
            local.get $2
            i32.add
            local.tee $7
            i32.const 4
            i32.add
            i32.load
            local.tee $8
            i32.const 16
            i32.shl
            local.get $5
            i32.const 16
            i32.shr_u
            i32.or
            i32.store
            local.get $3
            i32.const 4
            i32.add
            local.get $7
            i32.const 8
            i32.add
            i32.load
            local.tee $5
            i32.const 16
            i32.shl
            local.get $8
            i32.const 16
            i32.shr_u
            i32.or
            i32.store
            local.get $3
            i32.const 8
            i32.add
            local.get $7
            i32.const 12
            i32.add
            i32.load
            local.tee $8
            i32.const 16
            i32.shl
            local.get $5
            i32.const 16
            i32.shr_u
            i32.or
            i32.store
            local.get $3
            i32.const 12
            i32.add
            local.get $7
            i32.const 16
            i32.add
            i32.load
            local.tee $5
            i32.const 16
            i32.shl
            local.get $8
            i32.const 16
            i32.shr_u
            i32.or
            i32.store
            local.get $2
            i32.const 16
            i32.add
            local.set $2
            local.get $4
            i32.const -16
            i32.add
            local.tee $4
            i32.const 17
            i32.gt_u
            br_if $loop_2
          end ;; $loop_2
          local.get $6
          local.get $2
          i32.add
          local.set $3
          local.get $1
          local.get $2
          i32.add
          i32.const 2
          i32.add
          local.set $1
          br $block_7
        end ;; $block_8
        local.get $3
        local.get $1
        i32.load
        local.tee $5
        i32.store8
        local.get $4
        i32.const -1
        i32.add
        local.set $4
        local.get $3
        i32.const 1
        i32.add
        local.set $6
        i32.const 0
        local.set $2
        loop $loop_3
          local.get $6
          local.get $2
          i32.add
          local.tee $3
          local.get $1
          local.get $2
          i32.add
          local.tee $7
          i32.const 4
          i32.add
          i32.load
          local.tee $8
          i32.const 24
          i32.shl
          local.get $5
          i32.const 8
          i32.shr_u
          i32.or
          i32.store
          local.get $3
          i32.const 4
          i32.add
          local.get $7
          i32.const 8
          i32.add
          i32.load
          local.tee $5
          i32.const 24
          i32.shl
          local.get $8
          i32.const 8
          i32.shr_u
          i32.or
          i32.store
          local.get $3
          i32.const 8
          i32.add
          local.get $7
          i32.const 12
          i32.add
          i32.load
          local.tee $8
          i32.const 24
          i32.shl
          local.get $5
          i32.const 8
          i32.shr_u
          i32.or
          i32.store
          local.get $3
          i32.const 12
          i32.add
          local.get $7
          i32.const 16
          i32.add
          i32.load
          local.tee $5
          i32.const 24
          i32.shl
          local.get $8
          i32.const 8
          i32.shr_u
          i32.or
          i32.store
          local.get $2
          i32.const 16
          i32.add
          local.set $2
          local.get $4
          i32.const -16
          i32.add
          local.tee $4
          i32.const 18
          i32.gt_u
          br_if $loop_3
        end ;; $loop_3
        local.get $6
        local.get $2
        i32.add
        local.set $3
        local.get $1
        local.get $2
        i32.add
        i32.const 1
        i32.add
        local.set $1
      end ;; $block_7
      block $block_11
        local.get $4
        i32.const 16
        i32.and
        i32.eqz
        br_if $block_11
        local.get $3
        local.get $1
        i32.load16_u align=1
        i32.store16 align=1
        local.get $3
        local.get $1
        i32.load8_u offset=2
        i32.store8 offset=2
        local.get $3
        local.get $1
        i32.load8_u offset=3
        i32.store8 offset=3
        local.get $3
        local.get $1
        i32.load8_u offset=4
        i32.store8 offset=4
        local.get $3
        local.get $1
        i32.load8_u offset=5
        i32.store8 offset=5
        local.get $3
        local.get $1
        i32.load8_u offset=6
        i32.store8 offset=6
        local.get $3
        local.get $1
        i32.load8_u offset=7
        i32.store8 offset=7
        local.get $3
        local.get $1
        i32.load8_u offset=8
        i32.store8 offset=8
        local.get $3
        local.get $1
        i32.load8_u offset=9
        i32.store8 offset=9
        local.get $3
        local.get $1
        i32.load8_u offset=10
        i32.store8 offset=10
        local.get $3
        local.get $1
        i32.load8_u offset=11
        i32.store8 offset=11
        local.get $3
        local.get $1
        i32.load8_u offset=12
        i32.store8 offset=12
        local.get $3
        local.get $1
        i32.load8_u offset=13
        i32.store8 offset=13
        local.get $3
        local.get $1
        i32.load8_u offset=14
        i32.store8 offset=14
        local.get $3
        local.get $1
        i32.load8_u offset=15
        i32.store8 offset=15
        local.get $3
        i32.const 16
        i32.add
        local.set $3
        local.get $1
        i32.const 16
        i32.add
        local.set $1
      end ;; $block_11
      block $block_12
        local.get $4
        i32.const 8
        i32.and
        i32.eqz
        br_if $block_12
        local.get $3
        local.get $1
        i32.load8_u
        i32.store8
        local.get $3
        local.get $1
        i32.load8_u offset=1
        i32.store8 offset=1
        local.get $3
        local.get $1
        i32.load8_u offset=2
        i32.store8 offset=2
        local.get $3
        local.get $1
        i32.load8_u offset=3
        i32.store8 offset=3
        local.get $3
        local.get $1
        i32.load8_u offset=4
        i32.store8 offset=4
        local.get $3
        local.get $1
        i32.load8_u offset=5
        i32.store8 offset=5
        local.get $3
        local.get $1
        i32.load8_u offset=6
        i32.store8 offset=6
        local.get $3
        local.get $1
        i32.load8_u offset=7
        i32.store8 offset=7
        local.get $3
        i32.const 8
        i32.add
        local.set $3
        local.get $1
        i32.const 8
        i32.add
        local.set $1
      end ;; $block_12
      block $block_13
        local.get $4
        i32.const 4
        i32.and
        i32.eqz
        br_if $block_13
        local.get $3
        local.get $1
        i32.load8_u
        i32.store8
        local.get $3
        local.get $1
        i32.load8_u offset=1
        i32.store8 offset=1
        local.get $3
        local.get $1
        i32.load8_u offset=2
        i32.store8 offset=2
        local.get $3
        local.get $1
        i32.load8_u offset=3
        i32.store8 offset=3
        local.get $3
        i32.const 4
        i32.add
        local.set $3
        local.get $1
        i32.const 4
        i32.add
        local.set $1
      end ;; $block_13
      block $block_14
        local.get $4
        i32.const 2
        i32.and
        i32.eqz
        br_if $block_14
        local.get $3
        local.get $1
        i32.load8_u
        i32.store8
        local.get $3
        local.get $1
        i32.load8_u offset=1
        i32.store8 offset=1
        local.get $3
        i32.const 2
        i32.add
        local.set $3
        local.get $1
        i32.const 2
        i32.add
        local.set $1
      end ;; $block_14
      local.get $4
      i32.const 1
      i32.and
      i32.eqz
      br_if $block_1
      local.get $3
      local.get $1
      i32.load8_u
      i32.store8
    end ;; $block_1
    local.get $0
    )
  
  (func $145 (type $19)
    (param $0 i32)
    (param $1 i32)
    (param $2 i32)
    (result i32)
    (local $3 i32)
    (local $4 i32)
    (local $5 i32)
    (local $6 i64)
    block $block
      local.get $2
      i32.eqz
      br_if $block
      local.get $0
      local.get $1
      i32.store8
      local.get $2
      local.get $0
      i32.add
      local.tee $3
      i32.const -1
      i32.add
      local.get $1
      i32.store8
      local.get $2
      i32.const 3
      i32.lt_u
      br_if $block
      local.get $0
      local.get $1
      i32.store8 offset=2
      local.get $0
      local.get $1
      i32.store8 offset=1
      local.get $3
      i32.const -3
      i32.add
      local.get $1
      i32.store8
      local.get $3
      i32.const -2
      i32.add
      local.get $1
      i32.store8
      local.get $2
      i32.const 7
      i32.lt_u
      br_if $block
      local.get $0
      local.get $1
      i32.store8 offset=3
      local.get $3
      i32.const -4
      i32.add
      local.get $1
      i32.store8
      local.get $2
      i32.const 9
      i32.lt_u
      br_if $block
      local.get $0
      i32.const 0
      local.get $0
      i32.sub
      i32.const 3
      i32.and
      local.tee $4
      i32.add
      local.tee $3
      local.get $1
      i32.const 255
      i32.and
      i32.const 16843009
      i32.mul
      local.tee $1
      i32.store
      local.get $3
      local.get $2
      local.get $4
      i32.sub
      i32.const -4
      i32.and
      local.tee $4
      i32.add
      local.tee $2
      i32.const -4
      i32.add
      local.get $1
      i32.store
      local.get $4
      i32.const 9
      i32.lt_u
      br_if $block
      local.get $3
      local.get $1
      i32.store offset=8
      local.get $3
      local.get $1
      i32.store offset=4
      local.get $2
      i32.const -8
      i32.add
      local.get $1
      i32.store
      local.get $2
      i32.const -12
      i32.add
      local.get $1
      i32.store
      local.get $4
      i32.const 25
      i32.lt_u
      br_if $block
      local.get $3
      local.get $1
      i32.store offset=24
      local.get $3
      local.get $1
      i32.store offset=20
      local.get $3
      local.get $1
      i32.store offset=16
      local.get $3
      local.get $1
      i32.store offset=12
      local.get $2
      i32.const -16
      i32.add
      local.get $1
      i32.store
      local.get $2
      i32.const -20
      i32.add
      local.get $1
      i32.store
      local.get $2
      i32.const -24
      i32.add
      local.get $1
      i32.store
      local.get $2
      i32.const -28
      i32.add
      local.get $1
      i32.store
      local.get $4
      local.get $3
      i32.const 4
      i32.and
      i32.const 24
      i32.or
      local.tee $5
      i32.sub
      local.tee $2
      i32.const 32
      i32.lt_u
      br_if $block
      local.get $1
      i64.extend_i32_u
      local.tee $6
      i64.const 32
      i64.shl
      local.get $6
      i64.or
      local.set $6
      local.get $3
      local.get $5
      i32.add
      local.set $1
      loop $loop
        local.get $1
        local.get $6
        i64.store
        local.get $1
        i32.const 24
        i32.add
        local.get $6
        i64.store
        local.get $1
        i32.const 16
        i32.add
        local.get $6
        i64.store
        local.get $1
        i32.const 8
        i32.add
        local.get $6
        i64.store
        local.get $1
        i32.const 32
        i32.add
        local.set $1
        local.get $2
        i32.const -32
        i32.add
        local.tee $2
        i32.const 31
        i32.gt_u
        br_if $loop
      end ;; $loop
    end ;; $block
    local.get $0
    )
  
  (func $146 (type $6)
    (param $0 i32)
    (result i32)
    (local $1 i32)
    (local $2 i32)
    (local $3 i32)
    local.get $0
    local.set $1
    block $block
      block $block_0
        block $block_1
          local.get $0
          i32.const 3
          i32.and
          i32.eqz
          br_if $block_1
          block $block_2
            local.get $0
            i32.load8_u
            br_if $block_2
            local.get $0
            local.get $0
            i32.sub
            return
          end ;; $block_2
          local.get $0
          i32.const 1
          i32.add
          local.set $1
          loop $loop
            local.get $1
            i32.const 3
            i32.and
            i32.eqz
            br_if $block_1
            local.get $1
            i32.load8_u
            local.set $2
            local.get $1
            i32.const 1
            i32.add
            local.tee $3
            local.set $1
            local.get $2
            i32.eqz
            br_if $block_0
            br $loop
          end ;; $loop
        end ;; $block_1
        local.get $1
        i32.const -4
        i32.add
        local.set $1
        loop $loop_0
          local.get $1
          i32.const 4
          i32.add
          local.tee $1
          i32.load
          local.tee $2
          i32.const -1
          i32.xor
          local.get $2
          i32.const -16843009
          i32.add
          i32.and
          i32.const -2139062144
          i32.and
          i32.eqz
          br_if $loop_0
        end ;; $loop_0
        block $block_3
          local.get $2
          i32.const 255
          i32.and
          br_if $block_3
          local.get $1
          local.get $0
          i32.sub
          return
        end ;; $block_3
        loop $loop_1
          local.get $1
          i32.load8_u offset=1
          local.set $2
          local.get $1
          i32.const 1
          i32.add
          local.tee $3
          local.set $1
          local.get $2
          br_if $loop_1
          br $block
        end ;; $loop_1
      end ;; $block_0
      local.get $3
      i32.const -1
      i32.add
      local.set $3
    end ;; $block
    local.get $3
    local.get $0
    i32.sub
    )
  
  (data $26 (i32.const 65536)
    "\00\01\1c\02\1d\0e\18\03\1e\16\14\0f\19\11\04\08\1f\1b\0d\17\15\13\10\07\1a\0c\12\06\0b\05\n\09\00\01\02\02\03\03\03\03\04\04\04\04\04\04\04\04\05\05\05\05\05\05\05\05\05\05\05\05\05\05\05\05"
    "\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\06\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07"
    "\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\07\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08"
    "\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08"
    "\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\08\01\00\00\00\00\00\00\00S*\88\9c\dc\06\f3E\81\bbi\ea\a3\afN\f8\83 \87x2+\c5\b6"
    "\c3\8aw\bdW\c7\a2\fc\c0\a5\9b\84t\dcHn9$\a3\13\d4]\92\0c}NSb9\d6;\fc\09\89v]\eak\04\d1\82\fc\behM\e08\d3\e9s\a6\06s\ac#svP\f2{\17\bf\089"
    "\1f0\0b\bc\ff/\f1?>\da\14\b4\16#^\den\13\96O\9e\fdl\daZ\d4\bf<\cd~\8c\9e\f6\e2\cb\d7\8fuO\ea\d4%\a5\14\efs\0ez\10\ba\1a?b\bf\f6\d7W\d7\f6\f8\8d`\06\ac"
    "\c8\01\01\00\03\00\00\00nil\00\00\00\00\00\d8\01\01\002\00\00\00used block is not valid to be freed or r"
    "eallocated/proc/self/exenil pointer dereferencepanic: runtime er"
    "ror: index out of rangeslice out of rangeinvalid channel stateun"
    "reachablepanic: GC sweepsizertSize\00\00\00\00\00\00\b0\02\01\00\14\00\00\00allocation too l"
    "argeGC cycle\09live:\09\09\09\09\09live bytes:\09\09\09\09frees:\09\09\09\09\09allocs:\09\09\09\09\09fre"
    "ed bytes:\09\09\09sweep bytes:\09\09\09total bytes:\09\09\09last sweep:\09\09\09\09last sw"
    "eep bytes:\09\09last mark time:\09\09\09last graph time:\09\09last sweep time:"
    "\09\09last GC time:\09\09\c2\b5sgcInitHeapgcFreehi moontrade!\00\00\00\00\00\00\00\09\00\00\00\09\00\00\00")
  
  (data $27 (i32.const 66496)
    "\cc\03\01\00\01\00\00\00\01\00\00\00\n\02\01\00\0e\00\00\00\08\04\01\00\00\00\00\00\01\00\00\00\01\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\000\05\01\00")
  
  ;;(custom_section "producers"
  ;;  (after data)
  ;;  "\02\08language\01\03C99\00\0cprocessed-by\01\05c"
  ;;  "lang\\11.0.0 (https://github.com/"
  ;;  "tinygo-org/llvm-project 9ecb19f7"
  ;;  "74994a3efff5a6b89aa43ba2b8d2dd23"
  ;;  ")")
  
  )
#!/bin/bash

# Finder module for jterrazz command system
# This file defines all Finder visual formatting commands

# Main finder command handler
j_finder() {
    if [ $# -eq 0 ]; then
        j_finder_help
        return 1
    fi

    local subcommand="$1"
    shift

    case "$subcommand" in
        "format")
            j_finder_format "$@"
            ;;
        "reset")
            j_finder_reset "$@"
            ;;
        "help"|"-h"|"--help")
            j_finder_help
            ;;
        *)
            echo "❌ Unknown finder subcommand: $subcommand"
            j_finder_help
            return 1
            ;;
    esac
}

# Format folder with consistent visual settings (simple approach)
j_finder_format() {
    local folder_path=""
    local view_type="icon"
    local icon_size=64
    local grid_spacing=72
    local text_size=12
    local recursive=true

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --view)
                view_type="$2"
                shift 2
                ;;
            --icon-size)
                icon_size="$2"
                shift 2
                ;;
            --grid-spacing)
                grid_spacing="$2"
                shift 2
                ;;
            --text-size)
                text_size="$2"
                shift 2
                ;;
            --no-recursive)
                recursive=false
                shift
                ;;
            --nuclear)
                echo "💥 NUCLEAR OPTION: Resetting ALL Finder preferences"
                echo "⚠️  This will reset your entire Finder configuration. Continue? (y/N)"
                read -r response
                if [[ "$response" =~ ^[Yy]$ ]]; then
                    j_finder_nuclear_reset
                    return 0
                else
                    echo "❌ Nuclear reset cancelled"
                    return 0
                fi
                ;;
            --help|-h)
                j_finder_format_help
                return 0
                ;;
            -*)
                echo "❌ Unknown option: $1"
                j_finder_format_help
                return 1
                ;;
            *)
                if [ -z "$folder_path" ]; then
                    folder_path="$1"
                else
                    echo "❌ Multiple folder paths provided: $folder_path and $1"
                    return 1
                fi
                shift
                ;;
        esac
    done

    # Validate folder path
    if [ -z "$folder_path" ]; then
        echo "❌ Folder path is required"
        j_finder_format_help
        return 1
    fi

    # Expand tilde and resolve path
    folder_path="${folder_path/#\~/$HOME}"
    folder_path="$(realpath "$folder_path" 2>/dev/null)"

    if [ ! -d "$folder_path" ]; then
        echo "❌ Folder does not exist: $folder_path"
        return 1
    fi

    # Validate view type
    if [[ "$view_type" != "icon" && "$view_type" != "list" && "$view_type" != "column" && "$view_type" != "gallery" ]]; then
        echo "❌ Invalid view type: $view_type (must be: icon, list, column, gallery)"
        return 1
    fi

    echo "🎨 Formatting folder: $folder_path (simple mode - no permissions required)"
    echo "   Target view: $view_type, Icon Size: $icon_size, Grid: $grid_spacing, Text: $text_size"
    if [ "$recursive" = true ]; then
        echo "   Applying recursively to all subfolders..."
    fi

    # Step 1: Reset .DS_Store files to clear existing preferences
    echo "🗑️  Removing .DS_Store files to reset folder preferences..."
    if [ "$recursive" = true ]; then
        find "$folder_path" -name ".DS_Store" -delete 2>/dev/null
        local ds_count=$(find "$folder_path" -name ".DS_Store" 2>/dev/null | wc -l)
        echo "   Reset preferences for folder and all subfolders"
    else
        rm -f "$folder_path/.DS_Store" 2>/dev/null
        echo "   Reset preferences for target folder only"
    fi

    # Step 2: Set global Finder defaults (will apply to folders without .DS_Store)
    echo "🔧 Setting Finder global defaults..."
    
    # Set preferred view style
    local view_code=""
    case "$view_type" in
        "icon")   view_code="icnv" ;;
        "list")   view_code="Nlsv" ;;  
        "column") view_code="clmv" ;;
        "gallery") view_code="Flwv" ;;
    esac
    
    # Apply global Finder settings - SIMPLIFIED approach to avoid nesting errors
    defaults write com.apple.finder FXPreferredViewStyle -string "$view_code"
    defaults write com.apple.finder _FXSortFoldersFirst -bool true
    defaults write com.apple.finder _FXSortFoldersFirstOnDesktop -bool true
    
    # Force arrangement by name for icon views (most important setting)
    defaults write com.apple.finder FXArrangeGroupViewBy -string "name"
    defaults write com.apple.finder FXPreferredGroupBy -string "name"
    
    # Set specific view preferences using individual commands (no nesting)
    if [ "$view_type" = "icon" ]; then
        echo "   • Setting icon view preferences..."
        defaults write com.apple.finder FK_DefaultIconViewSettings -string '{"iconSize":'$icon_size',"gridSpacing":'$grid_spacing',"textSize":'$text_size',"arrangeBy":"name","labelOnBottom":true}'
        defaults write com.apple.finder DesktopViewOptions -string '{"iconSize":'$icon_size',"gridSpacing":'$grid_spacing',"textSize":'$text_size',"arrangeBy":"name","labelOnBottom":true}'
        
    elif [ "$view_type" = "list" ]; then
        echo "   • Setting list view preferences..."
        defaults write com.apple.finder FK_DefaultListViewSettings -string '{"sortColumn":"name","textSize":'$text_size',"reverseSort":false}'
        
    elif [ "$view_type" = "column" ]; then
        echo "   • Setting column view preferences..."  
        defaults write com.apple.finder FK_DefaultColumnViewSettings -string '{"textSize":'$text_size',"showPreviewColumn":true}'
    fi
    
    # Most important: Force "Arrange by Name" for ALL folder types
    # This is the key setting that ensures folders are always sorted by name
    defaults write com.apple.finder FXDefaultArrangement -string "name"
    defaults write com.apple.finder FXDefaultArrangeByName -bool true
    
    # Step 3: Restart Finder to apply changes
    echo "🔄 Restarting Finder to apply changes..."
    killall Finder 2>/dev/null
    
    # Wait a moment for Finder to restart
    sleep 2
    
    echo "✅ Folder formatting completed successfully"
    echo ""
    echo "📌 What was applied:"
    echo "   • Removed .DS_Store files (reset existing folder preferences)"
    echo "   • Set global Finder view to: $view_type"
    echo "   • Enabled 'folders first' sorting"  
    echo "   • FORCED 'Arrange by Name' as default for ALL folders"
    if [ "$view_type" = "icon" ]; then
        echo "   • Set icon size: ${icon_size}px, grid: ${grid_spacing}px, text: ${text_size}pt"
    fi
    echo "   • Restarted Finder to apply changes"
    echo ""
    echo "🔍 Verification:"
    echo "   Current view style: $(defaults read com.apple.finder FXPreferredViewStyle 2>/dev/null || echo 'not set')"
    echo "   Default arrangement: $(defaults read com.apple.finder FXDefaultArrangement 2>/dev/null || echo 'not set')"  
    echo "   Folders first: $(defaults read com.apple.finder _FXSortFoldersFirst 2>/dev/null || echo 'not set')"
    echo ""
    echo "💡 Open $folder_path in Finder to verify sorting by name is applied"
    echo ""
    echo "🛠️  If sorting STILL doesn't work, try these manual fixes:"
    echo "   1. Open the folder in Finder"
    echo "   2. View → Arrange By → Name"
    echo "   3. View → Show View Options (Cmd+J) → 'Use as Defaults'"
    echo ""
    echo "   Or run the NUCLEAR option:"
    echo "   j finder format $folder_path --nuclear"
}

# Reset folder to system defaults
j_finder_reset() {
    local folder_path=""
    local recursive=true

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --no-recursive)
                recursive=false
                shift
                ;;
            --help|-h)
                echo "Reset folder view settings to system defaults"
                echo ""
                echo "Usage: j finder reset <folder_path> [--no-recursive]"
                echo ""
                echo "Options:"
                echo "  --no-recursive    Only reset the target folder, not subfolders"
                return 0
                ;;
            -*)
                echo "❌ Unknown option: $1"
                return 1
                ;;
            *)
                if [ -z "$folder_path" ]; then
                    folder_path="$1"
                else
                    echo "❌ Multiple folder paths provided"
                    return 1
                fi
                shift
                ;;
        esac
    done

    if [ -z "$folder_path" ]; then
        echo "❌ Folder path is required"
        return 1
    fi

    # Expand and validate path
    folder_path="${folder_path/#\~/$HOME}"
    folder_path="$(realpath "$folder_path" 2>/dev/null)"

    if [ ! -d "$folder_path" ]; then
        echo "❌ Folder does not exist: $folder_path"
        return 1
    fi

    echo "🔄 Resetting folder preferences: $folder_path"
    
    if [ "$recursive" = true ]; then
        echo "🗑️  Removing all .DS_Store files recursively..."
        find "$folder_path" -name ".DS_Store" -delete 2>/dev/null
        echo "✅ Reset completed for all folders"
    else
        echo "🗑️  Removing .DS_Store file..."
        rm -f "$folder_path/.DS_Store" 2>/dev/null
        echo "✅ Reset completed for target folder"
    fi

    echo "💡 Refresh Finder windows to see default settings applied"
}

# Format help function
j_finder_format_help() {
    echo "🎨 Format Folder Visual Settings (Simple Mode)"
    echo ""
    echo "Usage: j finder format <folder_path> [options]"
    echo ""
    echo "Options:"
    echo "  --view <type>         View type: icon, list, column, gallery (default: icon)"
    echo "  --icon-size <size>    Icon size in pixels (default: 64, for icon view)"
    echo "  --grid-spacing <size> Grid spacing in pixels (default: 72, for icon view)"
    echo "  --text-size <size>    Text size in points (default: 12)"
    echo "  --no-recursive        Apply only to target folder, not subfolders"
    echo ""
    echo "How it works:"
    echo "  • Removes .DS_Store files to reset existing preferences"
    echo "  • Sets global Finder defaults for the chosen view type"
    echo "  • Ensures folders are sorted by name and folders-first"
    echo "  • Restarts Finder to apply changes immediately"
    echo ""
    echo "Examples:"
    echo "  j finder format ~/Documents                        # Icon view, 64px icons"
    echo "  j finder format ~/Downloads --view list            # List view, sorted by name"
    echo "  j finder format . --icon-size 128 --no-recursive   # Large icons, current folder only"
    echo ""
    echo "✅ No permissions required - works by setting global defaults + clearing local preferences"
}

# Main help function
j_finder_help() {
    echo "📁 Finder Visual Formatting Commands"
    echo ""
    echo "Usage: j finder <subcommand>"
    echo ""
    echo "Subcommands:"
    echo "  format    Apply consistent visual formatting to a folder"
    echo "  reset     Reset folder view settings to system defaults"
    echo "  help      Show this help"
    echo ""
    echo "Examples:"
    echo "  j finder format ~/Documents --view icon"
    echo "  j finder format ~/Downloads --view list"
    echo "  j finder reset ~/Desktop"
    echo ""
    echo "For detailed help:"
    echo "  j finder format --help"
    echo "  j finder reset --help"
}

# Auto-completion for finder subcommands
j_finder_completion() {
    echo "format reset help"
}

# Module metadata
J_MODULE_NAME="finder"
J_MODULE_DESCRIPTION="Format Finder folder views and settings"
J_MODULE_COMMANDS="format reset help"